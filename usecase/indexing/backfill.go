package indexing

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"

	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/pkg/errors"
)

var (
	ErrBackfillRunning = errors.New("backfill already running (use force flag to override it)")
)

type backfillUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client figmentclient.Client
}

func NewBackfillUseCase(cfg *config.Config, db *store.Store, c figmentclient.Client) *backfillUseCase {
	return &backfillUseCase{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

type BackfillUseCaseConfig struct {
	Parallel  bool
	Force     bool
	TargetIds []int64
}

func (uc *backfillUseCase) Execute(ctx context.Context, useCaseConfig BackfillUseCaseConfig) error {
	if err := uc.canExecute(); err != nil {
		return err
	}

	indexingPipeline, err := indexer.NewPipeline(uc.cfg, uc.db, uc.client)
	if err != nil {
		return err
	}

	return indexingPipeline.Backfill(ctx, indexer.BackfillConfig{
		Parallel:  useCaseConfig.Parallel,
		Force:     useCaseConfig.Force,
		TargetIds: useCaseConfig.TargetIds,
	})
}

// canExecute checks if reindex is already running
// if is it running we skip indexing
func (uc *backfillUseCase) canExecute() error {
	if _, err := uc.db.Reports.FindNotCompletedByKind(model.ReportKindSequentialReindex, model.ReportKindParallelReindex); err != nil {
		if err == store.ErrNotFound {
			return nil
		}
		return err
	}
	return ErrBackfillRunning
}
