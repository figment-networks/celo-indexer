package store

import (
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

func NewValidatorSeqStore(db *gorm.DB) *ValidatorSeqStore {
	return &ValidatorSeqStore{scoped(db, model.ValidatorSeq{})}
}

// ValidatorSeqStore handles operations on validators
type ValidatorSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the validator if it does not exist
func (s ValidatorSeqStore) CreateIfNotExists(validator *model.ValidatorSeq) error {
	_, err := s.FindByHeight(validator.Height)
	if isNotFound(err) {
		return s.Create(validator)
	}
	return nil
}

// FindByHeightAndAddress finds validator by height and address
func (s ValidatorSeqStore) FindByHeightAndAddress(height int64, address string) (*model.ValidatorSeq, error) {
	q := model.ValidatorSeq{
		Address: address,
	}
	var result model.ValidatorSeq

	err := s.db.
		Where(&q).
		Where("height = ?", height).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindByHeight finds validator sequences by height
func (s ValidatorSeqStore) FindByHeight(h int64) ([]model.ValidatorSeq, error) {
	var result []model.ValidatorSeq

	err := s.db.
		Where("height = ?", h).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindMostRecent finds most recent validator era sequence
func (s *ValidatorSeqStore) FindMostRecent() (*model.ValidatorSeq, error) {
	validatorSeq := &model.ValidatorSeq{}
	if err := findMostRecent(s.db, "time", validatorSeq); err != nil {
		return nil, err
	}
	return validatorSeq, nil
}

// FindLastByAddress finds last validator sequences for given address
func (s ValidatorSeqStore) FindLastByAddress(address string, limit int64) ([]model.ValidatorSeq, error) {
	q := model.ValidatorSeq{
		Address: address,
	}
	var result []model.ValidatorSeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// DeleteOlderThan deletes validator sequence older than given threshold
func (s *ValidatorSeqStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.ValidatorSeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

type ValidatorSeqSummary struct {
	Address    string         `json:"address"`
	TimeBucket types.Time     `json:"time_bucket"`
	SignedAvg  float64        `json:"signed_avg"`
	SignedMin  int64          `json:"signed_min"`
	SignedMax  int64          `json:"signed_max"`
	ScoreAvg   types.Quantity `json:"score_avg"`
	ScoreMin   types.Quantity `json:"score_min"`
	ScoreMax   types.Quantity `json:"score_max"`
}

// Summarize gets the summarized version of validator sequences
func (s *ValidatorSeqStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]ValidatorSeqSummary, error) {
	defer logQueryDuration(time.Now(), "ValidatorSeqStore_Summarize")

	tx := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Select(summarizeValidatorsForEraQuerySelect, interval).
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

	var models []ValidatorSeqSummary
	for rows.Next() {
		var summary ValidatorSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}
