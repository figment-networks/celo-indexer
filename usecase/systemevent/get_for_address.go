package systemevent

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getForAddressUseCase struct {
	db *psql.Store
}

func NewGetForAddressUseCase(db *psql.Store) *getForAddressUseCase {
	return &getForAddressUseCase{
		db: db,
	}
}

func (uc *getForAddressUseCase) Execute(address string, minHeight *int64, kind *model.SystemEventKind) (*ListView, error) {
	systemEvents, err := uc.db.GetCore().SystemEvents.FindByActor(address, store.FindSystemEventByActorQuery{
		Kind:      kind,
		MinHeight: minHeight,
	})
	if err != nil {
		return nil, err
	}

	return ToListView(systemEvents), nil
}
