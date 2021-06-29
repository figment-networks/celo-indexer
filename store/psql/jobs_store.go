package psql

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/jinzhu/gorm"
)

var _ store.Jobs = (*Jobs)(nil)

func NewJobsStore(db *gorm.DB) *Jobs {
	return &Jobs{scoped(db, model.Job{})}
}

// Jobs handles operations on jobs
type Jobs struct {
	baseStore
}

// Create inserts the job records
func (s Jobs) Create(jobs []model.Job) error {
	for _, job := range jobs {
		err := s.db.Create(&job).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// Update updates the selected fields of the job record
func (s Jobs) Update(job *model.Job, fields ...string) error {
	return s.db.Model(job).Select(fields).Updates(job).Error
}

// FindByHeight retrieves a job record for a given height
func (s Jobs) FindByHeight(height int64) (*model.Job, error) {
	var job model.Job

	err := s.db.Where("height = ?", height).Take(&job).Error
	if err != nil {
		return nil, err
	}

	return &job, err
}

// FindAllUnfinished retrieves all of the unfinished jobs
func (s Jobs) FindAllUnfinished() ([]model.Job, error) {
	var jobs []model.Job

	err := s.db.
		Where("finished_at IS NULL").
		Order("height ASC").
		Find(&jobs).
		Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// LastFinishedHeight returns the most recent finished height
func (s Jobs) LastFinishedHeight() (int64, error) {
	var result struct {
		Height int64
	}

	err := s.db.
		Table("jobs").
		Where("finished_at IS NOT NULL").
		Select("MAX(height) as height").
		Scan(&result).
		Error

	if err != nil || result.Height == 0 {
		return -1, err
	}

	return result.Height, nil
}