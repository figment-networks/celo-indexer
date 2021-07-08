package psql

import (
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"
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

// Create bulk-inserts the job records
func (s Jobs) Create(jobs []model.Job) error {
	now := time.Now()

	return bulk.Import(s.db, bulkInsertJobs, len(jobs), func(i int) bulk.Row {
		return bulk.Row{*jobs[i].Height, now, now}
	})
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
	var result struct{ Height int64 }

	err := s.db.
		Table("jobs").
		Where("finished_at IS NOT NULL").
		Select("MAX(height) as height").
		Scan(&result).
		Error

	if err != nil {
		return 0, err
	}

	return result.Height, nil
}

// LastSyncedHeight returns the most recent synced height
func (s Jobs) LastSyncedHeight() (int64, error) {
	var result struct{ Height int64 }

	err := s.db.
		Table("jobs").
		Where("finished_at IS NOT NULL").
		Where("height < (SELECT MIN(height) FROM jobs WHERE finished_at IS NULL)").
		Select("MAX(height) AS height").
		Scan(&result).
		Error

	if err != nil {
		return 0, err
	}

	return result.Height, nil
}
