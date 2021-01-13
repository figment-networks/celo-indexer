package psql

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

var _ store.Syncables = (*Syncables)(nil)

func NewSyncablesStore(db *gorm.DB) *Syncables {
	return &Syncables{scoped(db, model.Report{})}
}

// Syncables handles operations on syncables
type Syncables struct {
	baseStore
}

// Save creates the syncable
func (s Syncables) Save(val *model.Syncable) error {
	return s.baseStore.Save(val)
}

// FindSmallestIndexVersion returns smallest index version
func (s Syncables) FindSmallestIndexVersion() (*int64, error) {
	result := &model.Syncable{}

	err := s.db.
		Where("processed_at IS NOT NULL").
		Order("index_version").
		First(result).Error

	return &result.IndexVersion, checkErr(err)
}

// Exists returns true if a syncable exists at give height
func (s Syncables) FindByHeight(height int64) (syncable *model.Syncable, err error) {
	result := &model.Syncable{}

	err = s.db.
		Where("height = ?", height).
		First(result).
		Error

	return result, checkErr(err)
}

// FindMostRecent returns the most recent syncable
func (s Syncables) FindMostRecent() (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Order("height desc").
		First(result).Error

	return result, checkErr(err)
}

// FindLastInEpochForHeight finds last_in_epoch syncable for given height
func (s Syncables) FindLastInEpochForHeight(height int64) (syncable *model.Syncable, err error) {
	result := &model.Syncable{}

	err = s.db.
		Where("height >= ? AND last_in_epoch = ?", height, true).
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}

// FindLastInEpoch finds last syncable in given epoch
func (s Syncables) FindLastInEpoch(epoch int64) (syncable *model.Syncable, err error) {
	result := &model.Syncable{}

	err = s.db.
		Where("epoch = ?", epoch).
		Order("height DESC").
		First(result).
		Error

	return result, checkErr(err)
}

// FindFirstByDifferentIndexVersion returns first syncable with different index version
func (s Syncables) FindFirstByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Not("index_version = ?", indexVersion).
		Order("height").
		First(result).Error

	return result, checkErr(err)
}

// FindMostRecentByDifferentIndexVersion returns the most recent syncable with different index version
func (s Syncables) FindMostRecentByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Not("index_version = ?", indexVersion).
		Order("height desc").
		First(result).Error

	return result, checkErr(err)
}

// CreateOrUpdate creates a new syncable or updates an existing one
func (s Syncables) CreateOrUpdate(val *model.Syncable) error {
	existing, err := s.FindByHeight(val.Height)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}
	return s.Update(existing)
}

// CreateOrUpdate creates a new syncable or updates an existing one
func (s Syncables) SetProcessedAtForRange(reportID types.ID, startHeight int64, endHeight int64) error {
	err := s.db.
		Exec("UPDATE syncables SET report_id = ?, processed_at = NULL WHERE height >= ? AND height <= ?", reportID, startHeight, endHeight).
		Error

	return checkErr(err)
}
