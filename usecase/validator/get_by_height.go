package validator

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client figmentclient.Client
}

func NewGetByHeightUseCase(cfg *config.Config, db *store.Store, client figmentclient.Client) *getByHeightUseCase {
	return &getByHeightUseCase{
		cfg:    cfg,
		db:     db,
		client: client,
	}
}

func (uc *getByHeightUseCase) Execute(height *int64) (SeqListView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return SeqListView{}, err
	}
	lastH := mostRecentSynced.Height

	// Show last synced height, if not provided
	if height == nil {
		height = &lastH
	}

	if *height > lastH {
		return SeqListView{}, errors.New("height is not indexed yet")
	}

	validatorSeqs, err := uc.db.ValidatorSeq.FindByHeight(*height)
	if len(validatorSeqs) == 0 || err != nil {
		syncable, err := uc.db.Syncables.FindByHeight(*height)
		if err != nil {
			return SeqListView{}, err
		}

		indexingPipeline, err := indexer.NewPipeline(uc.cfg, uc.db, uc.client)
		if err != nil {
			return SeqListView{}, err
		}

		ctx := context.Background()
		payload, err := indexingPipeline.Run(ctx, indexer.RunConfig{
			Height:           syncable.Height,
			DesiredTargetIDs: []int64{indexer.TargetIndexValidatorSequences},
			Dry:              true,
		})
		if err != nil {
			return SeqListView{}, err
		}

		validatorSeqs = append(payload.NewValidatorSequences, payload.UpdatedValidatorSequences...)
	}

	return ToSeqListView(validatorSeqs), nil
}
