package indexing

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/pkg/errors"
)

var (
	ErrRunningSequentialReindex = errors.New("indexing skipped because sequential reindex hasn't finished yet")
)

type startUseCase struct {
	cfg    *config.Config
	db     *psql.Store
	client figmentclient.Client
}

func NewStartUseCase(cfg *config.Config, db *psql.Store, c figmentclient.Client) *startUseCase {
	return &startUseCase{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (uc *startUseCase) Execute(ctx context.Context, batchSize int64) error {
	if err := uc.canExecute(); err != nil {
		return err
	}

	indexingPipeline, err := indexer.NewPipeline(
		uc.cfg,
		uc.client,
		uc.db.GetCore().Syncables,
		uc.db.GetCore().Database,
		uc.db.GetCore().Reports,
		uc.db.GetBlocks().BlockSeq,
		uc.db.GetValidators().ValidatorSeq,
		uc.db.GetAccounts().AccountActivitySeq,
		uc.db.GetValidatorGroups().ValidatorGroupSeq,
		uc.db.GetValidators().ValidatorAgg,
		uc.db.GetValidatorGroups().ValidatorGroupAgg,
		uc.db.GetGovernance().ProposalAgg,
		uc.db.GetCore().SystemEvents,
		uc.db.GetGovernance().GovernanceActivitySeq,
	)
	if err != nil {
		return err
	}

	return indexingPipeline.Start(ctx, indexer.IndexConfig{
		BatchSize: batchSize,
	})
}

// canExecute checks if sequential reindex is already running
// if is it running we skip indexing
func (uc *startUseCase) canExecute() error {
	if _, err := uc.db.GetCore().Reports.FindNotCompletedByKind(model.ReportKindSequentialReindex); err != nil {
		if err == psql.ErrNotFound {
			return nil
		}
		return err
	}
	return ErrRunningSequentialReindex
}
