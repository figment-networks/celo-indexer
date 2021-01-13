package validator

import (
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/pkg/errors"
)

type getForMinHeightUseCase struct {
	db *psql.Store
}

func NewGetForMinHeightUseCase(db *psql.Store) *getForMinHeightUseCase {
	return &getForMinHeightUseCase{
		db: db,
	}
}

func (uc *getForMinHeightUseCase) Execute(height *int64) (*AggListView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.GetCore().Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}
	lastH := mostRecentSynced.Height

	// Show last synced height, if not provided
	if height == nil {
		height = &lastH
	}

	if *height > lastH {
		return nil, errors.New("height is not indexed yet")
	}

	validatorAggs, err := uc.db.GetValidators().ValidatorAgg.GetAllForHeightGreaterThan(*height)
	if err != nil {
		return nil, err
	}

	return ToAggListView(validatorAggs), nil
}


