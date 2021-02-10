package validator

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	cfg    *config.Config
	db     *psql.Store
	client figmentclient.ClientIface
}

func NewGetByHeightUseCase(cfg *config.Config, db *psql.Store, client figmentclient.ClientIface) *getByHeightUseCase {
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

	validatorSeqs, err := uc.getValidatorSeqs(height)
	if err != nil {
		return SeqListView{}, err
	}

	return ToSeqListView(validatorSeqs), nil
}

func (uc *getByHeightUseCase) getValidatorSeqs(height *int64) ([]model.ValidatorSeq, error) {
	return uc.db.GetValidators().ValidatorSeq.FindByHeight(*height)
}
