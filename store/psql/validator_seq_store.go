package psql

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

var _ store.ValidatorSeq = (*ValidatorSeq)(nil)

func NewValidatorSeqStore(db *gorm.DB) *ValidatorSeq {
	return &ValidatorSeq{scoped(db, model.ValidatorSeq{})}
}

// ValidatorSeq handles operations on validators
type ValidatorSeq struct {
	baseStore
}

// BulkUpsert insert validator sequences in bulk
func (s ValidatorSeq) BulkUpsert(records []model.ValidatorSeq) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertValidatorSeqs, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.Height,
				r.Time,
				r.Address,
				r.Affiliation,
				r.Signed,
				r.Score.String(),
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// FindByHeightAndAddress finds validator by height and address
func (s ValidatorSeq) FindByHeightAndAddress(height int64, address string) (*model.ValidatorSeq, error) {
	result := model.ValidatorSeq{}

	err := s.db.
		Model(&model.ValidatorSeq{}).
		Select(joinedAggregateSelect).
		Joins("LEFT JOIN validator_aggregates on validator_sequences.address = validator_aggregates.address").
		Where("validator_sequences.address = ?", address).
		Where("validator_sequences.height = ?", height).
		Limit(1).
		Scan(&result).Error

	return &result, checkErr(err)
}

// FindByHeight finds validator sequences by height
func (s ValidatorSeq) FindByHeight(height int64) ([]model.ValidatorSeq, error) {
	var result []model.ValidatorSeq

	err := s.db.
		Model(&model.ValidatorSeq{}).
		Select(joinedAggregateSelect).
		Joins("LEFT JOIN validator_aggregates on validator_sequences.address = validator_aggregates.address").
		Where("validator_sequences.height = ?", height).
		Scan(&result).Error

	return result, checkErr(err)
}

// FindMostRecent finds most recent validator era sequence
func (s *ValidatorSeq) FindMostRecent() (*model.ValidatorSeq, error) {
	result := model.ValidatorSeq{}

	err := s.db.
		Model(&model.ValidatorSeq{}).
		Select(joinedAggregateSelect).
		Joins("LEFT JOIN validator_aggregates on validator_sequences.address = validator_aggregates.address").
		Order("time DESC").
		Limit(1).
		Scan(&result).
		Error

	return &result, checkErr(err)
}

// FindLastByAddress finds last validator sequences for given address
func (s ValidatorSeq) FindLastByAddress(address string, limit int64) ([]model.ValidatorSeq, error) {
	var result []model.ValidatorSeq

	err := s.db.
		Model(&model.ValidatorSeq{}).
		Select(joinedAggregateSelect).
		Joins("LEFT JOIN validator_aggregates on validator_sequences.address = validator_aggregates.address").
		Where("validator_sequences.address = ?", address).
		Order("height DESC").
		Limit(limit).
		Scan(&result).
		Error

	return result, checkErr(err)
}

// DeleteOlderThan deletes validator sequence older than given threshold
func (s *ValidatorSeq) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.ValidatorSeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

// Summarize gets the summarized version of validator sequences
func (s *ValidatorSeq) Summarize(interval types.SummaryInterval, activityPeriods []store.ActivityPeriodRow) ([]store.ValidatorSeqSummary, error) {
	defer LogQueryDuration(time.Now(), "ValidatorSeqStore_Summarize")

	tx := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Select(summarizeValidatorsQuerySelect, interval).
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

	var models []store.ValidatorSeqSummary
	for rows.Next() {
		var summary store.ValidatorSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}
