package block

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
)

type getBlockSummaryUseCase struct {
	db *psql.Store
}

func NewGetBlockSummaryUseCase(db *psql.Store) *getBlockSummaryUseCase {
	return &getBlockSummaryUseCase{
		db: db,
	}
}

func (uc *getBlockSummaryUseCase) Execute(interval types.SummaryInterval, period string) ([]model.BlockSummary, error) {
	return uc.db.GetBlocks().BlockSummary.FindSummary(interval, period)
}
