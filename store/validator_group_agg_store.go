package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewValidatorGroupAggStore(db *gorm.DB) *ValidatorGroupAggStore {
	return &ValidatorGroupAggStore{scoped(db, model.ValidatorGroupAgg{})}
}

// ValidatorGroupAggStore handles operations on validators
type ValidatorGroupAggStore struct {
	baseStore
}

// CreateOrUpdate creates a new validator group or updates an existing one
func (s ValidatorGroupAggStore) CreateOrUpdate(val *model.ValidatorGroupAgg) error {
	existing, err := s.FindByAddress(val.Address)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}
	return s.Update(existing)
}

// FindBy returns an validator group for a matching attribute
func (s ValidatorGroupAggStore) FindBy(key string, value interface{}) (*model.ValidatorGroupAgg, error) {
	result := &model.ValidatorGroupAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns an validator group for the ID
func (s ValidatorGroupAggStore) FindByID(id int64) (*model.ValidatorGroupAgg, error) {
	return s.FindBy("id", id)
}

// FindByAddress return validator group by address
func (s *ValidatorGroupAggStore) FindByAddress(key string) (*model.ValidatorGroupAgg, error) {
	return s.FindBy("address", key)
}

// All returns all validator groups
func (s ValidatorGroupAggStore) All() ([]model.ValidatorGroupAgg, error) {
	var result []model.ValidatorGroupAgg

	err := s.db.
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}
