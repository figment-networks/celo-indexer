package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
	"time"
)

const (
	SyncerPersistorTaskName                = "SyncerPersistor"
	BlockSeqPersistorTaskName              = "BlockSeqPersistor"
	ValidatorSeqPersistorTaskName          = "ValidatorSeqPersistor"
	ValidatorGroupSeqPersistorTaskName     = "ValidatorGroupSeqPersistor"
	AccountActivitySeqPersistorTaskName    = "AccountActivitySeqPersistor"
	ValidatorAggPersistorTaskName          = "ValidatorAggPersistor"
	ValidatorGroupAggPersistorTaskName     = "ValidatorGroupAggPersistor"
	ProposalAggPersistorTaskName           = "ProposalAggPersistor"
	TaskNameSystemEventPersistor           = "SystemEventPersistor"
	GovernanceActivitySeqPersistorTaskName = "GovernanceActivitySeqPersistor"
)

// NewSyncerPersistorTask is responsible for storing syncable to persistence layer
func NewSyncerPersistorTask(syncableDb store.Syncables) pipeline.Task {
	return &syncerPersistorTask{
		syncableDb: syncableDb,
	}
}

type syncerPersistorTask struct {
	syncableDb store.Syncables
}

func (t *syncerPersistorTask) GetName() string {
	return SyncerPersistorTaskName
}

func (t *syncerPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return t.syncableDb.CreateOrUpdate(payload.Syncable)
}

// NewBlockSeqPersistorTask is responsible for storing block to persistence layer
func NewBlockSeqPersistorTask(blockSeqDb store.BlockSeq) pipeline.Task {
	return &blockSeqPersistorTask{
		blockSeqDb: blockSeqDb,
	}
}

type blockSeqPersistorTask struct {
	blockSeqDb store.BlockSeq
}

func (t *blockSeqPersistorTask) GetName() string {
	return BlockSeqPersistorTaskName
}

func (t *blockSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	if payload.NewBlockSequence != nil {
		return t.blockSeqDb.Create(payload.NewBlockSequence)
	}

	if payload.UpdatedBlockSequence != nil {
		return t.blockSeqDb.Save(payload.UpdatedBlockSequence)
	}

	return nil
}

// NewValidatorSeqPersistorTask is responsible for storing validator info to persistence layer
func NewValidatorSeqPersistorTask(validatorSeqDb store.ValidatorSeq) pipeline.Task {
	return &validatorSeqPersistorTask{
		validatorSeqDb: validatorSeqDb,
	}
}

type validatorSeqPersistorTask struct {
	validatorSeqDb store.ValidatorSeq
}

func (t *validatorSeqPersistorTask) GetName() string {
	return ValidatorSeqPersistorTaskName
}

func (t *validatorSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	if err := t.validatorSeqDb.BulkUpsert(payload.ValidatorSequences); err != nil {
		return err
	}

	return nil
}

// NewAccountActivitySeqPersistorTask is responsible for storing validator info to persistence layer
func NewAccountActivitySeqPersistorTask(accountActivitySeqDb store.AccountActivitySeq) pipeline.Task {
	return &accountActivitySeqPersistorTask{
		accountActivitySeqDb: accountActivitySeqDb,
	}
}

type accountActivitySeqPersistorTask struct {
	accountActivitySeqDb store.AccountActivitySeq
}

func (t *accountActivitySeqPersistorTask) GetName() string {
	return AccountActivitySeqPersistorTaskName
}

func (t *accountActivitySeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	// Delete current account activities first since it is not possible to find unique activities since the same activity kind could be found multiple times in one block
	_, err := t.accountActivitySeqDb.DeleteForHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	if err := t.accountActivitySeqDb.BulkUpsert(payload.AccountActivitySequences); err != nil {
		return err
	}

	return nil
}

// NewValidatorGroupSeqPersistorTask is responsible for storing validator era info to persistence layer
func NewValidatorGroupSeqPersistorTask(validatorGroupSeqDb store.ValidatorGroupSeq) pipeline.Task {
	return &validatorGroupSeqPersistorTask{
		validatorGroupSeqDb: validatorGroupSeqDb,
	}
}

type validatorGroupSeqPersistorTask struct {
	validatorGroupSeqDb store.ValidatorGroupSeq
}

func (t *validatorGroupSeqPersistorTask) GetName() string {
	return ValidatorGroupSeqPersistorTaskName
}

func (t *validatorGroupSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	if err := t.validatorGroupSeqDb.BulkUpsert(payload.ValidatorGroupSequences); err != nil {
		return err
	}

	return nil
}

// NewValidatorAggPersistorTask store validator aggregate to persistence layer
func NewValidatorAggPersistorTask(validatorAggDb store.ValidatorAgg) pipeline.Task {
	return &validatorAggPersistorTask{
		validatorAggDb: validatorAggDb,
	}
}

type validatorAggPersistorTask struct {
	validatorAggDb store.ValidatorAgg
}

func (t *validatorAggPersistorTask) GetName() string {
	return ValidatorAggPersistorTaskName
}

func (t *validatorAggPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, aggregate := range payload.NewValidatorAggregates {
		if err := t.validatorAggDb.Create(&aggregate); err != nil {
			return err
		}
	}

	for _, aggregate := range payload.UpdatedValidatorAggregates {
		if err := t.validatorAggDb.Save(&aggregate); err != nil {
			return err
		}
	}

	return nil
}

// NewValidatorGroupAggPersistorTask psql validator group aggregate to persistence layer
func NewValidatorGroupAggPersistorTask(validatorGroupAggDb store.ValidatorGroupAgg) pipeline.Task {
	return &validatorGroupAggPersistorTask{
		validatorGroupAggDb: validatorGroupAggDb,
	}
}

type validatorGroupAggPersistorTask struct {
	validatorGroupAggDb store.ValidatorGroupAgg
}

func (t *validatorGroupAggPersistorTask) GetName() string {
	return ValidatorGroupAggPersistorTaskName
}

func (t *validatorGroupAggPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, aggregate := range payload.NewValidatorGroupAggregates {
		if err := t.validatorGroupAggDb.Create(&aggregate); err != nil {
			return err
		}
	}

	for _, aggregate := range payload.UpdatedValidatorGroupAggregates {
		if err := t.validatorGroupAggDb.Save(&aggregate); err != nil {
			return err
		}
		if err := t.validatorGroupAggDb.CalculateCumulativeUptime(aggregate.Address); err != nil {
			return err
		}
	}

	return nil
}

// NewProposalAggPersistorTask psql validator aggregate to persistence layer
func NewProposalAggPersistorTask(proposalAggDb store.ProposalAgg) pipeline.Task {
	return &proposalAggPersistorTask{
		proposalAggDb: proposalAggDb,
	}
}

type proposalAggPersistorTask struct {
	proposalAggDb store.ProposalAgg
}

func (t *proposalAggPersistorTask) GetName() string {
	return ProposalAggPersistorTaskName
}

func (t *proposalAggPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, aggregate := range payload.NewProposalAggregates {
		if err := t.proposalAggDb.Create(&aggregate); err != nil {
			return err
		}
	}

	for _, aggregate := range payload.UpdatedProposalAggregates {
		if err := t.proposalAggDb.Save(&aggregate); err != nil {
			return err
		}
	}

	return nil
}

//NewSystemEventPersistorTask psql system events to persistance layer
func NewSystemEventPersistorTask(systemEventDb store.SystemEvents) pipeline.Task {
	return &systemEventPersistorTask{
		systemEventDb:  systemEventDb,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameSystemEventPersistor),
	}
}

type systemEventPersistorTask struct {
	systemEventDb  store.SystemEvents
	metricObserver metrics.Observer
}

func (t *systemEventPersistorTask) GetName() string {
	return TaskNameSystemEventPersistor
}

func (t *systemEventPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	if err := t.systemEventDb.BulkUpsert(payload.SystemEvents); err != nil {
		return err
	}

	return nil
}

// NewGovernanceActivitySeqPersistorTask is responsible for storing validator info to persistence layer
func NewGovernanceActivitySeqPersistorTask(governanceActivitySeqDb store.GovernanceActivitySeq) pipeline.Task {
	return &governanceActivitySeqPersistorTask{
		governanceActivitySeqDb: governanceActivitySeqDb,
	}
}

type governanceActivitySeqPersistorTask struct {
	governanceActivitySeqDb store.GovernanceActivitySeq
}

func (t *governanceActivitySeqPersistorTask) GetName() string {
	return GovernanceActivitySeqPersistorTaskName
}

func (t *governanceActivitySeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	// Delete current governance activities first
	_, err := t.governanceActivitySeqDb.DeleteForHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	if err := t.governanceActivitySeqDb.BulkUpsert(payload.GovernanceActivitySequences); err != nil {
		return err
	}

	return nil
}
