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

var _ store.ValidatorGroupSummary = (*ValidatorGroupSummary)(nil)

func NewValidatorGroupSummaryStore(db *gorm.DB) *ValidatorGroupSummary {
	return &ValidatorGroupSummary{scoped(db, model.ValidatorGroupSummary{})}
}

// ValidatorGroupSummary handles operations on validator group summaries
type ValidatorGroupSummary struct {
	baseStore
}

// BulkUpsert insert validator sequences in bulk
func (s ValidatorGroupSummary) BulkUpsert(records []model.ValidatorGroupSummary) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertValidatorGroupSummaries, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.TimeInterval,
				r.TimeBucket,
				r.IndexVersion,
				r.Address,
				r.CommissionAvg.String(),
				r.CommissionMax.String(),
				r.CommissionMin.String(),
				r.ActiveVotesAvg.String(),
				r.ActiveVotesMax.String(),
				r.ActiveVotesMin.String(),
				r.PendingVotesAvg.String(),
				r.PendingVotesMax.String(),
				r.PendingVotesMin.String(),
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Find find validator group summary by query
func (s ValidatorGroupSummary) Find(query *model.ValidatorGroupSummary) (*model.ValidatorGroupSummary, error) {
	var result model.ValidatorGroupSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindActivityPeriods Finds activity periods
func (s *ValidatorGroupSummary) FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]store.ActivityPeriodRow, error) {
	defer metrics.LogQueryDuration(time.Now(), "ValidatorGroupSummaryStore_FindActivityPeriods")

	rows, err := s.db.
		Raw(validatorGroupSummaryActivityPeriodsQuery, fmt.Sprintf("1%s", interval), interval, indexVersion).
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

// FindSummary gets summary for validator group summary
func (s *ValidatorGroupSummary) FindSummary(interval types.SummaryInterval, period string) ([]store.ValidatorGroupSummaryRow, error) {
	defer metrics.LogQueryDuration(time.Now(), "ValidatorGroupSummaryStore_FindSummary")

	rows, err := s.db.
		Raw(allValidatorGroupsSummaryForIntervalQuery, interval, period, interval).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []store.ValidatorGroupSummaryRow
	for rows.Next() {
		var row store.ValidatorGroupSummaryRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// FindSummaryByAddress gets summary for given validator group
func (s *ValidatorGroupSummary) FindSummaryByAddress(address string, interval types.SummaryInterval, period string) ([]model.ValidatorGroupSummary, error) {
	defer metrics.LogQueryDuration(time.Now(), "ValidatorGroupSummaryStore_FindSummaryByAddress")

	rows, err := s.db.Raw(validatorGroupSummaryForIntervalQuery, interval, period, address, interval).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.ValidatorGroupSummary
	for rows.Next() {
		var row model.ValidatorGroupSummary
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// FindMostRecent finds most recent validator group summary
func (s *ValidatorGroupSummary) FindMostRecent() (*model.ValidatorGroupSummary, error) {
	validatorSummary := &model.ValidatorGroupSummary{}
	err := findMostRecent(s.db, "time_bucket", validatorSummary)
	return validatorSummary, checkErr(err)
}

// FindMostRecentByInterval finds most recent validator group summary for interval
func (s *ValidatorGroupSummary) FindMostRecentByInterval(interval types.SummaryInterval) (*model.ValidatorGroupSummary, error) {
	query := &model.ValidatorGroupSummary{
		Summary: &model.Summary{TimeInterval: interval},
	}
	result := model.ValidatorGroupSummary{}

	err := s.db.
		Where(query).
		Order("time_bucket DESC").
		Take(&result).
		Error

	return &result, checkErr(err)
}

// DeleteOlderThan deleted validator group summary records older than given threshold
func (s *ValidatorGroupSummary) DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error) {
	statement := s.db.
		Unscoped().
		Where("time_interval = ? AND time_bucket < ?", interval, purgeThreshold).
		Delete(&model.ValidatorGroupSummary{})

	if statement.Error != nil {
		return nil, checkErr(statement.Error)
	}

	return &statement.RowsAffected, nil
}
