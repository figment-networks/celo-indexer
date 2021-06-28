package validatorgroup

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	cfg    *config.Config
	db     *psql.Store
	client figmentclient.Client
}

func NewGetByHeightUseCase(cfg *config.Config, db *psql.Store, client figmentclient.Client) *getByHeightUseCase {
	return &getByHeightUseCase{
		cfg:    cfg,
		db:     db,
		client: client,
	}
}

func (uc *getByHeightUseCase) Execute(height *int64) (SeqListView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.GetCore().Syncables.FindMostRecent()
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

	validatorSeqs, err := uc.db.GetValidatorGroups().ValidatorGroupSeq.FindByHeight(*height)
	if len(validatorSeqs) == 0 || err != nil {
		syncable, err := uc.db.GetCore().Syncables.FindByHeight(*height)
		if err != nil {
			return SeqListView{}, err
		}

		indexingPipeline, err := indexer.NewPipeline(
			uc.cfg,
			uc.client,
			nil,
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
			return SeqListView{}, err
		}

		ctx := context.Background()
		payload, err := indexingPipeline.Run(ctx, indexer.RunConfig{
			Height:           syncable.Height,
			DesiredTargetIDs: []int64{indexer.TargetIndexValidatorGroupSequences},
			Dry:              true,
		})
		if err != nil {
			return SeqListView{}, err
		}

		validatorSeqs = payload.ValidatorGroupSequences
	}

	return ToSeqListView(validatorSeqs), nil
}
