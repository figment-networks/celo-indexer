package psql

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"
	"github.com/jinzhu/gorm"
	"time"
)

var _ store.SystemEvents = (*SystemEvents)(nil)

func NewSystemEventsStore(db *gorm.DB) *SystemEvents {
	return &SystemEvents{scoped(db, model.SystemEvent{})}
}

// SystemEvents handles operations on syncables
type SystemEvents struct {
	baseStore
}

// BulkUpsert insert system events in bulk
func (s SystemEvents) BulkUpsert(records []model.SystemEvent) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertSystemEvents, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.Height,
				r.Time,
				r.Actor,
				r.Kind,
				r.Data,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// FindByHeight returns system events by height
func (s SystemEvents) FindByHeight(height int64) ([]model.SystemEvent, error) {
	var result []model.SystemEvent

	err := s.db.
		Where("height = ?", height).
		Find(result).
		Error

	return result, checkErr(err)
}

// FindByActor returns system events by actor
func (s SystemEvents) FindByActor(actorAddress string, query store.FindSystemEventByActorQuery) ([]model.SystemEvent, error) {
	var result []model.SystemEvent
	q := model.SystemEvent{}
	if query.Kind != nil {
		q.Kind = *query.Kind
	}

	statement := s.db.
		Where("actor = ?", actorAddress).
		Where(&q)

	if query.MinHeight != nil {
		statement = statement.Where("height > ?", query.MinHeight)
	}

	err := statement.
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s SystemEvents) FindAll(query store.FindAll) ([]model.SystemEvent, error) {
	var result []model.SystemEvent

	statement := s.db
	offset := (query.Page - 1) * query.Limit
	statement = statement.Offset(offset).Limit(query.Limit)

	err := statement.
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindUnique returns unique system
func (s SystemEvents) FindUnique(height int64, address string, kind model.SystemEventKind) (*model.SystemEvent, error) {
	q := model.SystemEvent{
		Height: height,
		Actor:  address,
		Kind:   kind,
	}

	var result model.SystemEvent
	err := s.db.
		Where(&q).
		First(&result).
		Error

	return &result, checkErr(err)
}

// CreateOrUpdate creates a new system event or updates an existing one
func (s SystemEvents) CreateOrUpdate(val *model.SystemEvent) error {
	existing, err := s.FindUnique(val.Height, val.Actor, val.Kind)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}

	existing.Update(*val)

	return s.Save(existing)
}

// FindMostRecent finds most recent system event
func (s *SystemEvents) FindMostRecent() (*model.SystemEvent, error) {
	systemEvent := &model.SystemEvent{}
	if err := findMostRecent(s.db, "time", systemEvent); err != nil {
		return nil, err
	}
	return systemEvent, nil
}

// DeleteOlderThan deletes system events older than given threshold
func (s *SystemEvents) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.BlockSeq{})

	return &tx.RowsAffected, checkErr(tx.Error)
}
