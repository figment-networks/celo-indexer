package validatorgroup

import (
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getByAddressUseCase struct {
	db *psql.Store
}

func NewGetByAddressUseCase(db *psql.Store) *getByAddressUseCase {
	return &getByAddressUseCase{
		db: db,
	}
}

func (uc *getByAddressUseCase) Execute(address string, sequencesLimit int64) (*AggDetailsView, error) {
	validatorGroupAggs, err := uc.db.GetValidatorGroups().ValidatorGroupAgg.FindByAddress(address)
	if err != nil {
		return nil, err
	}

	sequences, err := uc.getSessionSequences(address, sequencesLimit)
	if err != nil {
		return nil, err
	}

	delegations, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, indexer.OperationTypeValidatorGroupVoteActivatedReceived, sequencesLimit)

	return ToAggDetailsView(validatorGroupAggs, sequences, delegations), nil
}

func (uc *getByAddressUseCase) getSessionSequences(address string, sequencesLimit int64) ([]model.ValidatorGroupSeq, error) {
	var sequences []model.ValidatorGroupSeq
	var err error
	if sequencesLimit > 0 {
		sequences, err = uc.db.GetValidatorGroups().ValidatorGroupSeq.FindLastByAddress(address, sequencesLimit)
		if err != nil {
			return nil, err
		}
	}
	return sequences, nil
}
