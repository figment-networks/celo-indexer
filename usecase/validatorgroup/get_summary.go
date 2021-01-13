package validatorgroup

import (
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
)

type getSummaryUseCase struct {
	db *psql.Store
}

func NewGetSummaryUseCase(db *psql.Store) *getSummaryUseCase {
	return &getSummaryUseCase{
		db: db,
	}
}

func (uc *getSummaryUseCase) Execute(interval types.SummaryInterval, period string, address string) (interface{}, error) {
	if address == "" {
		return uc.db.GetValidatorGroups().ValidatorGroupSummary.FindSummary(interval, period)
	}
	return uc.db.GetValidatorGroups().ValidatorGroupSummary.FindSummaryByAddress(address, interval, period)
}

