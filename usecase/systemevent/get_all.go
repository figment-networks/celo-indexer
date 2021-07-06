package systemevent

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getAllUseCase struct {
	db *psql.Store
}

func NewGetAllUseCase(db *psql.Store) *getAllUseCase {
	return &getAllUseCase{
		db: db,
	}
}

func (uc *getAllUseCase) Execute(page int64, limit int64) (*ListView, error) {
	systemEvents, err := uc.db.GetCore().SystemEvents.FindAll(store.FindAll{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	return ToListView(systemEvents), nil
}
