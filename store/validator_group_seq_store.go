package store

import (
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

func NewValidatorGroupSeqStore(db *gorm.DB) *ValidatorGroupSeqStore {
	return &ValidatorGroupSeqStore{scoped(db, model.ValidatorGroupSeq{})}
}

// ValidatorGroupSeqStore handles operations on validator groups
type ValidatorGroupSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the validator group if it does not exist
func (s ValidatorGroupSeqStore) CreateIfNotExists(validator *model.ValidatorGroupSeq) error {
	_, err := s.FindByHeight(validator.Height)
	if isNotFound(err) {
		return s.Create(validator)
	}
	return nil
}

// FindByHeightAndAddress finds validator group by height and address
func (s ValidatorGroupSeqStore) FindByHeightAndAddress(height int64, address string) (*model.ValidatorGroupSeq, error) {
	q := model.ValidatorGroupSeq{
		Address: address,
	}
	var result model.ValidatorGroupSeq

	err := s.db.
		Where(&q).
		Where("height = ?", height).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindByHeight finds validator group sequences by height
func (s ValidatorGroupSeqStore) FindByHeight(h int64) ([]model.ValidatorGroupSeq, error) {
	var result []model.ValidatorGroupSeq

	err := s.db.
		Where("height = ?", h).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindLastByAddress finds last validator group sequences for given address
func (s ValidatorGroupSeqStore) FindLastByAddress(address string, limit int64) ([]model.ValidatorGroupSeq, error) {
	q := model.ValidatorGroupSeq{
		Address: address,
	}
	var result []model.ValidatorGroupSeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindMostRecent finds most recent validator group sequence
func (s *ValidatorGroupSeqStore) FindMostRecent() (*model.ValidatorGroupSeq, error) {
	validatorSeq := &model.ValidatorGroupSeq{}
	if err := findMostRecent(s.db, "time", validatorSeq); err != nil {
		return nil, err
	}
	return validatorSeq, nil
}

// DeleteOlderThan deletes validator group sequence older than given threshold
func (s *ValidatorGroupSeqStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.ValidatorGroupSeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

type ValidatorGroupSeqSummary struct {
	Address            string         `json:"address"`
	TimeBucket         types.Time     `json:"time_bucket"`
	CommissionAvg      types.Quantity `json:"commission_avg"`
	CommissionMin      types.Quantity `json:"commission_min"`
	CommissionMax      types.Quantity `json:"commission_max"`
	ActiveVotesAvg     types.Quantity `json:"active_votes_avg"`
	ActiveVotesMin     types.Quantity `json:"active_votes_min"`
	ActiveVotesMax     types.Quantity `json:"active_votes_max"`
	ActiveVoteUnitsAvg types.Quantity `json:"active_vote_units_avg"`
	ActiveVoteUnitsMin types.Quantity `json:"active_vote_units_min"`
	ActiveVoteUnitsMax types.Quantity `json:"active_vote_units_max"`
	PendingVotesAvg    types.Quantity `json:"pending_votes_avg"`
	PendingVotesMin    types.Quantity `json:"pending_votes_min"`
	PendingVotesMax    types.Quantity `json:"pending_votes_max"`
}

// Summarize gets the summarized version of validator sequences
func (s *ValidatorGroupSeqStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]ValidatorGroupSeqSummary, error) {
	defer logQueryDuration(time.Now(), "ValidatorGroupSeqStore_Summarize")

	tx := s.db.
		Table(model.ValidatorGroupSeq{}.TableName()).
		Select(summarizeValidatorsForSessionQuerySelect, interval).
		Order("time_bucket").
		Group("address, time_bucket")

	if len(activityPeriods) == 1 {
		activityPeriod := activityPeriods[0]
		tx = tx.Or("time < ? OR time >= ?", activityPeriod.Min, activityPeriod.Max)
	} else {
		for i, activityPeriod := range activityPeriods {
			isLast := i == len(activityPeriods)-1

			if isLast {
				tx = tx.Or("time >= ?", activityPeriod.Max)
			} else {
				duration, err := interval.Duration()
				if err != nil {
					return nil, err
				}
				tx = tx.Or("time >= ? AND time < ?", activityPeriod.Max.Add(*duration), activityPeriods[i+1].Min)
			}
		}
	}

	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ValidatorGroupSeqSummary
	for rows.Next() {
		var summary ValidatorGroupSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}
