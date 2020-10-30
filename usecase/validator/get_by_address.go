package validator

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
)

type getByAddressUseCase struct {
	db *store.Store
}

func NewGetByAddressUseCase(db *store.Store) *getByAddressUseCase {
	return &getByAddressUseCase{
		db: db,
	}
}

func (uc *getByAddressUseCase) Execute(address string, sequencesLimit int64) (*AggDetailsView, error) {
	validatorAggs, err := uc.db.ValidatorAgg.FindByAddress(address)
	if err != nil {
		return nil, err
	}

	sequences, err := uc.getSessionSequences(address, sequencesLimit)
	if err != nil {
		return nil, err
	}

	return ToAggDetailsView(validatorAggs, sequences), nil
}

func (uc *getByAddressUseCase) getSessionSequences(address string, sequencesLimit int64) ([]model.ValidatorSeq, error) {
	var sequences []model.ValidatorSeq
	var err error
	if sequencesLimit > 0 {
		sequences, err = uc.db.ValidatorSeq.FindLastByAddress(address, sequencesLimit)
		if err != nil {
			return nil, err
		}
	}
	return sequences, nil
}
