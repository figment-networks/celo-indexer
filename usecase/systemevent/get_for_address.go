package systemevent

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
)

type getForAddressUseCase struct {
	db *store.Store
}

func NewGetForAddressUseCase(db *store.Store) *getForAddressUseCase {
	return &getForAddressUseCase{
		db: db,
	}
}

func (uc *getForAddressUseCase) Execute(address string, minHeight *int64, kind *model.SystemEventKind) (*ListView, error) {
	systemEvents, err := uc.db.SystemEvents.FindByActor(address, store.FindSystemEventByActorQuery{
		Kind:      kind,
		MinHeight: minHeight,
	})
	if err != nil {
		return nil, err
	}

	return ToListView(systemEvents), nil
}
