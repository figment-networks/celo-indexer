package block

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getBlockTimesUseCase struct {
	db *psql.Store
}

func NewGetBlockTimesUseCase(db *psql.Store) *getBlockTimesUseCase {
	return &getBlockTimesUseCase{
		db: db,
	}
}

func (uc *getBlockTimesUseCase) Execute(limit int64) store.GetAvgRecentTimesResult {
	return uc.db.GetBlocks().BlockSeq.GetAvgRecentTimes(limit)
}

