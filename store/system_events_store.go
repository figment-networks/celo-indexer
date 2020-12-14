package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
	"time"
)

func NewSystemEventsStore(db *gorm.DB) *SystemEventsStore {
	return &SystemEventsStore{scoped(db, model.SystemEvent{})}
}

// SystemEventsStore handles operations on syncables
type SystemEventsStore struct {
	baseStore
}

// FindByHeight returns system events by height
func (s SystemEventsStore) FindByHeight(height int64) ([]model.SystemEvent, error) {
	var result []model.SystemEvent

	err := s.db.
		Where("height = ?", height).
		Find(result).
		Error

	return result, checkErr(err)
}

type FindSystemEventByActorQuery struct {
	Kind      *model.SystemEventKind
	MinHeight *int64
}

// FindByActor returns system events by actor
func (s SystemEventsStore) FindByActor(actorAddress string, query FindSystemEventByActorQuery) ([]model.SystemEvent, error) {
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

// FindUnique returns unique system
func (s SystemEventsStore) FindUnique(height int64, address string, kind model.SystemEventKind) (*model.SystemEvent, error) {
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
func (s SystemEventsStore) CreateOrUpdate(val *model.SystemEvent) error {
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
func (s *SystemEventsStore) FindMostRecent() (*model.SystemEvent, error) {
	systemEvent := &model.SystemEvent{}
	if err := findMostRecent(s.db, "time", systemEvent); err != nil {
		return nil, err
	}
	return systemEvent, nil
}

// DeleteOlderThan deletes system events older than given threshold
func (s *SystemEventsStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.BlockSeq{})

	return &tx.RowsAffected, checkErr(tx.Error)
}
