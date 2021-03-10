package psql

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
)

var _ store.ValidatorGroupSeq = (*ValidatorGroupSeq)(nil)

func NewValidatorGroupSeqStore(db *gorm.DB) *ValidatorGroupSeq {
	return &ValidatorGroupSeq{scoped(db, model.ValidatorGroupSeq{})}
}

// ValidatorGroupSeq handles operations on validator groups
type ValidatorGroupSeq struct {
	baseStore
}

// BulkUpsert insert validator group sequences in bulk
func (s ValidatorGroupSeq) BulkUpsert(records []model.ValidatorGroupSeq) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertValidatorGroupSeqs, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.Height,
				r.Time,
				r.Address,
				r.Commission.String(),
				r.ActiveVotes.String(),
				r.PendingVotes.String(),
				r.VotingCap.String(),
				r.MembersCount,
				r.MembersAvgSigned,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// FindByHeightAndAddress finds validator group by height and address
func (s ValidatorGroupSeq) FindByHeightAndAddress(height int64, address string) (*model.ValidatorGroupSeq, error) {
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
func (s ValidatorGroupSeq) FindByHeight(h int64) ([]model.ValidatorGroupSeq, error) {

	rows, err := s.db.Raw(validatorGroupByHeight, h).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.ValidatorGroupSeq
	for rows.Next() {
		var row model.ValidatorGroupSeq
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil

}

// FindLastByAddress finds last validator group sequences for given address
func (s ValidatorGroupSeq) FindLastByAddress(address string, limit int64) ([]model.ValidatorGroupSeq, error) {
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
func (s *ValidatorGroupSeq) FindMostRecent() (*model.ValidatorGroupSeq, error) {
	validatorSeq := &model.ValidatorGroupSeq{}
	if err := findMostRecent(s.db, "time", validatorSeq); err != nil {
		return nil, err
	}
	return validatorSeq, nil
}

// DeleteOlderThan deletes validator group sequence older than given threshold
func (s *ValidatorGroupSeq) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.ValidatorGroupSeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

// Summarize gets the summarized version of validator sequences
func (s *ValidatorGroupSeq) Summarize(interval types.SummaryInterval, activityPeriods []store.ActivityPeriodRow) ([]store.ValidatorGroupSeqSummary, error) {
	defer LogQueryDuration(time.Now(), "ValidatorGroupSeqStore_Summarize")

	tx := s.db.
		Table(model.ValidatorGroupSeq{}.TableName()).
		Select(summarizeValidatorGroupsQuerySelect, interval).
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

	var models []store.ValidatorGroupSeqSummary
	for rows.Next() {
		var summary store.ValidatorGroupSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}
