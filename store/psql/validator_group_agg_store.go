package psql

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/jinzhu/gorm"
)

var _ store.ValidatorGroupAgg = (*ValidatorGroupAgg)(nil)

func NewValidatorGroupAggStore(db *gorm.DB) *ValidatorGroupAgg {
	return &ValidatorGroupAgg{scoped(db, model.ValidatorGroupAgg{})}
}

// ValidatorGroupAgg handles operations on validators
type ValidatorGroupAgg struct {
	baseStore
}

// Create creates the validator group aggregate
func (s ValidatorGroupAgg) Create(val *model.ValidatorGroupAgg) error {
	return s.baseStore.Create(val)
}

// Save creates the validator group aggregate
func (s ValidatorGroupAgg) Save(val *model.ValidatorGroupAgg) error {
	return s.baseStore.Save(val)
}

// CreateOrUpdate creates a new validator group or updates an existing one
func (s ValidatorGroupAgg) CreateOrUpdate(val *model.ValidatorGroupAgg) error {
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
func (s ValidatorGroupAgg) FindBy(key string, value interface{}) (*model.ValidatorGroupAgg, error) {
	result := &model.ValidatorGroupAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns an validator group for the ID
func (s ValidatorGroupAgg) FindByID(id int64) (*model.ValidatorGroupAgg, error) {
	return s.FindBy("id", id)
}

// FindByAddress return validator group by address
func (s *ValidatorGroupAgg) FindByAddress(key string) (*model.ValidatorGroupAgg, error) {
	return s.FindBy("address", key)
}

// All returns all validator groups
func (s ValidatorGroupAgg) All() ([]model.ValidatorGroupAgg, error) {
	var result []model.ValidatorGroupAgg

	err := s.db.
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}

// CalculateCumulativeUptime calculate cumulative uptime avg by address
func (s *ValidatorGroupAgg) CalculateCumulativeUptime(address string) error {
	err := s.db.
		Exec(UpdateCalculate, address, address).
		Error

	return checkErr(err)
}
