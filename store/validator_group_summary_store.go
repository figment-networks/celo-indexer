package store

import (
	"fmt"
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

func NewValidatorGroupSummaryStore(db *gorm.DB) *ValidatorGroupSummaryStore {
	return &ValidatorGroupSummaryStore{scoped(db, model.ValidatorGroupSummary{})}
}

// ValidatorGroupSummaryStore handles operations on validator group summaries
type ValidatorGroupSummaryStore struct {
	baseStore
}

// Find find validator group summary by query
func (s ValidatorGroupSummaryStore) Find(query *model.ValidatorGroupSummary) (*model.ValidatorGroupSummary, error) {
	var result model.ValidatorGroupSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindActivityPeriods Finds activity periods
func (s *ValidatorGroupSummaryStore) FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error) {
	defer logQueryDuration(time.Now(), "ValidatorGroupSummaryStore_FindActivityPeriods")

	rows, err := s.db.
		Raw(validatorGroupSummaryActivityPeriodsQuery, fmt.Sprintf("1%s", interval), interval, indexVersion).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ActivityPeriodRow
	for rows.Next() {
		var row ActivityPeriodRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

type ValidatorGroupSummaryRow struct {
	TimeBucket         string         `json:"time_bucket"`
	TimeInterval       string         `json:"time_interval"`
	CommissionAvg      string `json:"commission_avg"`
	CommissionMin      string `json:"commission_min"`
	CommissionMax      string `json:"commission_max"`
	ActiveVotesAvg     string `json:"active_votes_avg"`
	ActiveVotesMin     string `json:"active_votes_min"`
	ActiveVotesMax     string `json:"active_votes_max"`
	ActiveVoteUnitsAvg string `json:"active_vote_units_avg"`
	ActiveVoteUnitsMin string `json:"active_vote_units_min"`
	ActiveVoteUnitsMax string `json:"active_vote_units_max"`
	PendingVotesAvg    string `json:"pending_votes_avg"`
	PendingVotesMin    string `json:"pending_votes_min"`
	PendingVotesMax    string `json:"pending_votes_max"`
}

// FindSummary gets summary for validator group summary
func (s *ValidatorGroupSummaryStore) FindSummary(interval types.SummaryInterval, period string) ([]ValidatorGroupSummaryRow, error) {
	defer logQueryDuration(time.Now(), "ValidatorGroupSummaryStore_FindSummary")

	rows, err := s.db.
		Raw(allValidatorGroupsSummaryForIntervalQuery, interval, period, interval).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ValidatorGroupSummaryRow
	for rows.Next() {
		var row ValidatorGroupSummaryRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// FindSummaryByAddress gets summary for given validator group
func (s *ValidatorGroupSummaryStore) FindSummaryByAddress(address string, interval types.SummaryInterval, period string) ([]model.ValidatorGroupSummary, error) {
	defer logQueryDuration(time.Now(), "ValidatorGroupSummaryStore_FindSummaryByAddress")

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
func (s *ValidatorGroupSummaryStore) FindMostRecent() (*model.ValidatorGroupSummary, error) {
	validatorSummary := &model.ValidatorGroupSummary{}
	err := findMostRecent(s.db, "time_bucket", validatorSummary)
	return validatorSummary, checkErr(err)
}

// FindMostRecentByInterval finds most recent validator group summary for interval
func (s *ValidatorGroupSummaryStore) FindMostRecentByInterval(interval types.SummaryInterval) (*model.ValidatorGroupSummary, error) {
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
func (s *ValidatorGroupSummaryStore) DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error) {
	statement := s.db.
		Unscoped().
		Where("time_interval = ? AND time_bucket < ?", interval, purgeThreshold).
		Delete(&model.ValidatorGroupSummary{})

	if statement.Error != nil {
		return nil, checkErr(statement.Error)
	}

	return &statement.RowsAffected, nil
}
