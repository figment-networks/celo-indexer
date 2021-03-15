package psql

import (
	"fmt"
	"time"

	"github.com/figment-networks/celo-indexer/metrics"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

var _ store.ValidatorSummary = (*ValidatorSummary)(nil)

func NewValidatorSummaryStore(db *gorm.DB) *ValidatorSummary {
	return &ValidatorSummary{scoped(db, model.ValidatorSummary{})}
}

// ValidatorSummary handles operations on validators
type ValidatorSummary struct {
	baseStore
}

// BulkUpsert insert validator sequences in bulk
func (s ValidatorSummary) BulkUpsert(records []model.ValidatorSummary) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertValidatorSummaries, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.TimeInterval,
				r.TimeBucket,
				r.IndexVersion,
				r.Address,
				r.ScoreAvg.String(),
				r.ScoreMax.String(),
				r.ScoreMin.String(),
				r.SignedAvg,
				r.SignedMax,
				r.SignedMin,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Find find validator summary by query
func (s ValidatorSummary) Find(query *model.ValidatorSummary) (*model.ValidatorSummary, error) {
	var result model.ValidatorSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindActivityPeriods Finds activity periods
func (s *ValidatorSummary) FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]store.ActivityPeriodRow, error) {
	defer metrics.LogQueryDuration(time.Now(), "ValidatorSummaryStore_FindActivityPeriods")

	rows, err := s.db.
		Raw(validatorSummaryActivityPeriodsQuery, fmt.Sprintf("1%s", interval), interval, indexVersion).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []store.ActivityPeriodRow
	for rows.Next() {
		var row store.ActivityPeriodRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// FindSummary gets summary for validator summary
func (s *ValidatorSummary) FindSummary(interval types.SummaryInterval, period string) ([]store.ValidatorSummaryRow, error) {
	defer metrics.LogQueryDuration(time.Now(), "ValidatorSummaryStore_FindSummary")

	rows, err := s.db.
		Raw(allValidatorsSummaryForIntervalQuery, interval, period, interval).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []store.ValidatorSummaryRow
	for rows.Next() {
		var row store.ValidatorSummaryRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// FindSummaryByAddress gets summary for given validator
func (s *ValidatorSummary) FindSummaryByAddress(address string, interval types.SummaryInterval, period string) ([]model.ValidatorSummary, error) {
	defer metrics.LogQueryDuration(time.Now(), "ValidatorSummaryStore_FindSummaryByAddress")

	rows, err := s.db.Raw(validatorSummaryForIntervalQuery, interval, period, address, interval).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.ValidatorSummary
	for rows.Next() {
		var row model.ValidatorSummary
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// FindMostRecent finds most recent validator summary
func (s *ValidatorSummary) FindMostRecent() (*model.ValidatorSummary, error) {
	validatorSummary := &model.ValidatorSummary{}
	err := findMostRecent(s.db, "time_bucket", validatorSummary)
	return validatorSummary, checkErr(err)
}

// FindMostRecentByInterval finds most recent validator summary for interval
func (s *ValidatorSummary) FindMostRecentByInterval(interval types.SummaryInterval) (*model.ValidatorSummary, error) {
	query := &model.ValidatorSummary{
		Summary: &model.Summary{TimeInterval: interval},
	}
	result := model.ValidatorSummary{}

	err := s.db.
		Where(query).
		Order("time_bucket DESC").
		Take(&result).
		Error

	return &result, checkErr(err)
}

// DeleteOlderThan deleted validator summary records older than given threshold
func (s *ValidatorSummary) DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error) {
	statement := s.db.
		Unscoped().
		Where("time_interval = ? AND time_bucket < ?", interval, purgeThreshold).
		Delete(&model.ValidatorSummary{})

	if statement.Error != nil {
		return nil, checkErr(statement.Error)
	}

	return &statement.RowsAffected, nil
}
