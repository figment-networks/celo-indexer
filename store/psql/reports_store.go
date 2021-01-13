package psql

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/jinzhu/gorm"
)

var _ store.Reports = (*Reports)(nil)

func NewReportsStore(db *gorm.DB) *Reports {
	return &Reports{scoped(db, model.Report{})}
}

// Reports handles operations on reports
type Reports struct {
	baseStore
}

// Create creates the report
func (s Reports) Create(val *model.Report) error {
	return s.baseStore.Create(val)
}

// Save creates the report
func (s Reports) Save(val *model.Report) error {
	return s.baseStore.Save(val)
}

// FindNotCompletedByIndexVersion returns the report by index version and kind
func (s Reports) FindNotCompletedByIndexVersion(indexVersion int64, kinds ...model.ReportKind) (*model.Report, error) {
	query := &model.Report{
		IndexVersion: indexVersion,
	}
	result := &model.Report{}

	err := s.db.
		Where(query).
		Where("kind IN(?)", kinds).
		Where("completed_at IS NULL").
		First(result).Error

	return result, checkErr(err)
}

// Last returns the last report
func (s Reports) FindNotCompletedByKind(kinds ...model.ReportKind) (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Where("kind IN(?)", kinds).
		Where("completed_at IS NULL").
		First(result).Error

	return result, checkErr(err)
}

// Last returns the last report
func (s Reports) Last() (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Order("id DESC").
		First(result).Error

	return result, checkErr(err)
}

// DeleteByKinds deletes reports with kind reindexing sequential or parallel
func (s *Reports) DeleteByKinds(kinds []model.ReportKind) error {
	err := s.db.
		Unscoped().
		Where("kind IN(?)", kinds).
		Delete(&model.Report{}).
		Error

	return checkErr(err)
}
