package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/datalake"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/pkg/errors"
)

const (
	CtxReport  = "context_report"
	maxRetries = 3
)

var (
	StageAnalyzer pipeline.StageName = "stage_analyzer"

	ErrIsPristine          = errors.New("cannot run because database is empty")
	ErrIndexCannotBeRun    = errors.New("cannot run index process")
	ErrBackfillCannotBeRun = errors.New("cannot run backfill process")
)

type indexingPipeline struct {
	cfg    *config.Config
	client figmentclient.Client

	syncableDb              store.Syncables
	jobsDb                  store.Jobs
	databaseDb              store.Database
	reportDb                store.Reports
	blockSeqDb              store.BlockSeq
	validatorSeqDb          store.ValidatorSeq
	accountActivitySeqDb    store.AccountActivitySeq
	validatorGroupSeqDb     store.ValidatorGroupSeq
	validatorAggDb          store.ValidatorAgg
	validatorGroupAggDb     store.ValidatorGroupAgg
	proposalAggDb           store.ProposalAgg
	systemEventDb           store.SystemEvents
	governanceActivitySeqDb store.GovernanceActivitySeq

	status       *pipelineStatus
	configParser ConfigParser
	pipeline     pipeline.CustomPipeline
}

func NewPipeline(
	cfg *config.Config,
	client figmentclient.Client,
	dl *datalake.DataLake,

	syncableDb store.Syncables,
	jobsDb store.Jobs,
	databaseDb store.Database,
	reportsDb store.Reports,
	blockSeqDb store.BlockSeq,
	validatorSeqDb store.ValidatorSeq,
	accountActivitySeqDb store.AccountActivitySeq,
	validatorGroupSeqDb store.ValidatorGroupSeq,
	validatorAggDb store.ValidatorAgg,
	validatorGroupAggDb store.ValidatorGroupAgg,
	proposalAggDb store.ProposalAgg,
	systemEventDb store.SystemEvents,
	governanceActivitySeqDb store.GovernanceActivitySeq,
) (*indexingPipeline, error) {
	p := pipeline.NewCustom(NewPayloadFactory(dl))

	// Setup logger
	p.SetLogger(NewLogger())

	// Setup stage
	p.AddStage(
		pipeline.NewStageWithTasks(pipeline.StageSetup, pipeline.RetryingTask(NewHeightMetaRetrieverTask(client), isTransient, maxRetries)),
	)

	// Fetcher stage
	if dl != nil {
		p.AddStage(
			pipeline.NewAsyncStageWithTasks(pipeline.StageFetcher,
				pipeline.RetryingTask(NewBlockDownloaderTask(), isTransient, maxRetries),
				pipeline.RetryingTask(NewValidatorDownloaderTask(), isTransient, maxRetries),
				pipeline.RetryingTask(NewValidatorGroupDownloaderTask(), isTransient, maxRetries),
				pipeline.RetryingTask(NewTransactionDownloaderTask(), isTransient, maxRetries),
			),
		)
	} else {
		p.AddStage(
			pipeline.NewAsyncStageWithTasks(pipeline.StageFetcher,
				pipeline.RetryingTask(NewBlockFetcherTask(client), isTransient, maxRetries),
				pipeline.RetryingTask(NewValidatorFetcherTask(client), isTransient, maxRetries),
				pipeline.RetryingTask(NewValidatorGroupFetcherTask(client), isTransient, maxRetries),
				pipeline.RetryingTask(NewTransactionFetcherTask(client), isTransient, maxRetries),
			),
		)
	}

	// Parser stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(pipeline.StageParser,
			pipeline.RetryingTask(NewGovernanceLogsParserTask(), isTransient, maxRetries),
		),
	)

	// Syncer stage
	p.AddStage(
		pipeline.NewStageWithTasks(pipeline.StageSyncer, pipeline.RetryingTask(NewMainSyncerTask(syncableDb), isTransient, maxRetries)),
	)

	// Sequencer stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(
			pipeline.StageSequencer,
			pipeline.RetryingTask(NewBlockSeqCreatorTask(blockSeqDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorSeqCreatorTask(cfg), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorGroupSeqCreatorTask(cfg), isTransient, maxRetries),
			pipeline.RetryingTask(NewAccountActivitySeqCreatorTask(cfg), isTransient, maxRetries),
			pipeline.RetryingTask(NewGovernanceActivitySeqCreatorTask(cfg), isTransient, maxRetries),
		),
	)

	// Aggregator stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(
			pipeline.StageAggregator,
			pipeline.RetryingTask(NewValidatorAggCreatorTask(client, validatorAggDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorGroupAggCreatorTask(client, validatorGroupAggDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewProposalAggCreatorTask(proposalAggDb), isTransient, maxRetries),
		),
	)

	// Analyzer stage
	p.AddStage(
		pipeline.NewStageWithTasks(
			StageAnalyzer,
			pipeline.RetryingTask(NewSystemEventCreatorTask(cfg, validatorSeqDb, accountActivitySeqDb, validatorGroupSeqDb), isTransient, maxRetries),
		),
	)

	// Persistor stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(
			pipeline.StagePersistor,
			pipeline.RetryingTask(NewSyncerPersistorTask(syncableDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewSystemEventPersistorTask(systemEventDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewBlockSeqPersistorTask(blockSeqDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorSeqPersistorTask(validatorSeqDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorGroupSeqPersistorTask(validatorGroupSeqDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewAccountActivitySeqPersistorTask(accountActivitySeqDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewGovernanceActivitySeqPersistorTask(governanceActivitySeqDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorAggPersistorTask(validatorAggDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorGroupAggPersistorTask(validatorGroupAggDb), isTransient, maxRetries),
			pipeline.RetryingTask(NewProposalAggPersistorTask(proposalAggDb), isTransient, maxRetries),
		),
	)

	// Create config parser
	configParser, err := NewConfigParser(cfg.IndexerConfigFile)
	if err != nil {
		return nil, err
	}

	statusChecker := pipelineStatusChecker{syncableDb, configParser.GetCurrentVersionId()}
	pipelineStatus, err := statusChecker.getStatus()
	if err != nil {
		return nil, err
	}

	return &indexingPipeline{
		cfg:    cfg,
		client: client,

		syncableDb:              syncableDb,
		jobsDb:                  jobsDb,
		databaseDb:              databaseDb,
		reportDb:                reportsDb,
		blockSeqDb:              blockSeqDb,
		validatorSeqDb:          validatorSeqDb,
		accountActivitySeqDb:    accountActivitySeqDb,
		validatorGroupSeqDb:     validatorGroupSeqDb,
		validatorAggDb:          validatorAggDb,
		validatorGroupAggDb:     validatorGroupAggDb,
		proposalAggDb:           proposalAggDb,
		systemEventDb:           systemEventDb,
		governanceActivitySeqDb: governanceActivitySeqDb,

		pipeline:     p,
		status:       pipelineStatus,
		configParser: configParser,
	}, nil
}

type IndexConfig struct {
	BatchSize   int64
	StartHeight int64
}

// Start starts indexing process
func (p *indexingPipeline) Start(ctx context.Context, indexCfg IndexConfig) error {
	if err := p.canRunIndex(); err != nil {
		return err
	}

	indexVersion := p.configParser.GetCurrentVersionId()

	source, err := NewIndexSource(p.cfg, p.client, p.syncableDb, p.jobsDb, &IndexSourceConfig{
		BatchSize:   indexCfg.BatchSize,
		StartHeight: indexCfg.StartHeight,
	})
	if err != nil {
		return err
	}

	sink := NewSink(p.syncableDb, p.databaseDb, p.client, indexVersion)

	reportCreator := &reportCreator{
		kind:         model.ReportKindIndex,
		indexVersion: indexVersion,
		startHeight:  source.startHeight,
		endHeight:    source.endHeight,
		reportDb:     p.reportDb,
	}

	if err := reportCreator.create(); err != nil {
		return err
	}

	versionIds := p.configParser.GetAllVersionedVersionIds()
	pipelineOptionsCreator := &pipelineOptionsCreator{
		configParser:      p.configParser,
		desiredVersionIds: versionIds,
	}
	pipelineOptions, err := pipelineOptionsCreator.parse()
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline [start=%d] [end=%d]", source.startHeight, source.endHeight))

	ctxWithReport := context.WithValue(ctx, CtxReport, reportCreator.report)
	err = p.pipeline.Start(ctxWithReport, source, sink, pipelineOptions)

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = reportCreator.complete(source.Len(), sink.successCount, err)
	return err
}

func (p *indexingPipeline) canRunIndex() error {
	if !p.status.isPristine && !p.status.isUpToDate {
		if p.configParser.IsAnyVersionSequential(p.status.missingVersionIds) {
			return ErrIndexCannotBeRun
		}
	}
	return nil
}

type BackfillConfig struct {
	Parallel  bool
	Force     bool
	TargetIds []int64
}

func (p *indexingPipeline) Backfill(ctx context.Context, backfillCfg BackfillConfig) error {
	if err := p.canRunBackfill(backfillCfg.Parallel); err != nil {
		return err
	}

	indexVersion := p.configParser.GetCurrentVersionId()

	source, err := NewBackfillSource(p.cfg, p.client, p.syncableDb, indexVersion)
	if err != nil {
		return err
	}

	sink := NewSink(p.syncableDb, p.databaseDb, p.client, indexVersion)

	kind := model.ReportKindSequentialReindex
	if backfillCfg.Parallel {
		kind = model.ReportKindParallelReindex
	}

	if backfillCfg.Force {
		if err := p.reportDb.DeleteByKinds([]model.ReportKind{model.ReportKindParallelReindex, model.ReportKindSequentialReindex}); err != nil {
			return err
		}
	}

	reportCreator := &reportCreator{
		kind:         kind,
		indexVersion: indexVersion,
		startHeight:  source.startHeight,
		endHeight:    source.endHeight,
		reportDb:     p.reportDb,
	}

	versionIds := p.status.missingVersionIds
	pipelineOptionsCreator := &pipelineOptionsCreator{
		configParser:      p.configParser,
		desiredVersionIds: versionIds,
	}
	pipelineOptions, err := pipelineOptionsCreator.parse()
	if err != nil {
		return err
	}

	if err := reportCreator.createIfNotExists(model.ReportKindSequentialReindex, model.ReportKindParallelReindex); err != nil {
		return err
	}

	if err := p.syncableDb.SetProcessedAtForRange(reportCreator.report.ID, source.startHeight, source.endHeight); err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, reportCreator.report)

	logger.Info(fmt.Sprintf("starting pipeline backfill [start=%d] [end=%d] [kind=%s]", source.startHeight, source.endHeight, kind))

	if err := p.pipeline.Start(ctxWithReport, source, sink, pipelineOptions); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = reportCreator.complete(source.Len(), sink.successCount, err)

	return err
}

func (p *indexingPipeline) canRunBackfill(isParallel bool) error {
	if p.status.isPristine {
		return ErrIsPristine
	}

	if !p.status.isUpToDate {
		if isParallel && p.configParser.IsAnyVersionSequential(p.status.missingVersionIds) {
			return ErrBackfillCannotBeRun
		}
	}
	return nil
}

type RunConfig struct {
	Height            int64
	DesiredVersionIDs []int64
	DesiredTargetIDs  []int64
	Dry               bool
}

func (p *indexingPipeline) Run(ctx context.Context, runCfg RunConfig) (*payload, error) {
	pipelineOptionsCreator := &pipelineOptionsCreator{
		configParser:      p.configParser,
		dry:               runCfg.Dry,
		desiredVersionIds: runCfg.DesiredVersionIDs,
		desiredTargetIds:  runCfg.DesiredTargetIDs,
	}

	pipelineOptions, err := pipelineOptionsCreator.parse()
	if err != nil {
		return nil, err
	}

	logger.Info("running pipeline...",
		logger.Field("height", runCfg.Height),
		logger.Field("versions", runCfg.DesiredVersionIDs),
		logger.Field("targets", runCfg.DesiredTargetIDs),
	)

	runPayload, err := p.pipeline.Run(ctx, runCfg.Height, pipelineOptions)
	if err != nil {
		logger.Info(fmt.Sprintf("pipeline completed with error [Err: %+v]", err))
		return nil, err
	}

	logger.Info("pipeline completed successfully")

	payload := runPayload.(*payload)
	return payload, nil
}

func isTransient(error) bool {
	return true
}

func RunFetcherPipeline(height int64, client figmentclient.Client, dl *datalake.DataLake) error {
	p := pipeline.NewCustom(NewPayloadFactory(dl))

	// Setup stage
	p.AddStage(
		pipeline.NewStageWithTasks(pipeline.StageSetup,
			pipeline.RetryingTask(NewHeightMetaRetrieverTask(client), isTransient, maxRetries),
		),
	)

	// Fetcher stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(pipeline.StageFetcher,
			pipeline.RetryingTask(NewBlockFetcherTask(client), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorFetcherTask(client), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorGroupFetcherTask(client), isTransient, maxRetries),
			pipeline.RetryingTask(NewTransactionFetcherTask(client), isTransient, maxRetries),
		),
	)

	// Persistor stage
	p.AddStage(
		pipeline.NewAsyncStageWithTasks(pipeline.StagePersistor,
			pipeline.RetryingTask(NewBlockUploaderTask(), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorUploaderTask(), isTransient, maxRetries),
			pipeline.RetryingTask(NewValidatorGroupUploaderTask(), isTransient, maxRetries),
			pipeline.RetryingTask(NewTransactionUploaderTask(), isTransient, maxRetries),
		),
	)

	_, err := p.Run(context.Background(), height, nil)
	if err != nil {
		return err
	}

	return nil
}
