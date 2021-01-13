package psql

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/jinzhu/gorm"
)

var _ store.ValidatorAgg = (*ValidatorAgg)(nil)

func NewValidatorAggStore(db *gorm.DB) *ValidatorAgg {
	return &ValidatorAgg{scoped(db, model.ValidatorAgg{})}
}

// ValidatorAgg handles operations on validators
type ValidatorAgg struct {
	baseStore
}

// Create creates the validator aggregate
func (s ValidatorAgg) Create(val *model.ValidatorAgg) error {
	return s.baseStore.Create(val)
}

// Save creates the validator aggregate
func (s ValidatorAgg) Save(val *model.ValidatorAgg) error {
	return s.baseStore.Save(val)
}

// FindBy returns an validator for a matching attribute
func (s ValidatorAgg) FindBy(key string, value interface{}) (*model.ValidatorAgg, error) {
	result := &model.ValidatorAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns an validator for the ID
func (s ValidatorAgg) FindByID(id int64) (*model.ValidatorAgg, error) {
	return s.FindBy("id", id)
}

// FindByAddress return validator by address
func (s *ValidatorAgg) FindByAddress(key string) (*model.ValidatorAgg, error) {
	return s.FindBy("address", key)
}

// GetAllForHeightGreaterThan returns validators who have been validating since given height
func (s *ValidatorAgg) GetAllForHeightGreaterThan(height int64) ([]model.ValidatorAgg, error) {
	var result []model.ValidatorAgg

	err := s.baseStore.db.
		Where("recent_as_validator_height >= ?", height).
		Find(&result).
		Error

	return result, checkErr(err)
}

// All returns all validators
func (s ValidatorAgg) All() ([]model.ValidatorAgg, error) {
	var result []model.ValidatorAgg

	err := s.db.
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}
