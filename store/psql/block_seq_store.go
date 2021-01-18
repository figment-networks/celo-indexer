package psql

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/figment-networks/celo-indexer/model"
)

var _ store.BlockSeq = (*BlockSeq)(nil)

func NewBlockSeqStore(db *gorm.DB) *BlockSeq {
	return &BlockSeq{scoped(db, model.BlockSeq{})}
}

// BlockSeq handles operations on blocks
type BlockSeq struct {
	baseStore
}

// Create creates the block
func (s BlockSeq) Create(block *model.BlockSeq) error {
	return s.baseStore.Create(block)
}

// Save saves the block
func (s BlockSeq) Save(block *model.BlockSeq) error {
	return s.baseStore.Save(block)
}

// CreateIfNotExists creates the block if it does not exist
func (s BlockSeq) CreateIfNotExists(block *model.BlockSeq) error {
	_, err := s.FindByHeight(block.Height)
	if isNotFound(err) {
		return s.Create(block)
	}
	return nil
}

// FindBy returns a block for a matching attribute
func (s BlockSeq) FindBy(key string, value interface{}) (*model.BlockSeq, error) {
	result := &model.BlockSeq{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a block with matching ID
func (s BlockSeq) FindByID(id int64) (*model.BlockSeq, error) {
	return s.FindBy("id", id)
}

// FindByHeight returns a block with the matching height
func (s BlockSeq) FindByHeight(height int64) (*model.BlockSeq, error) {
	return s.FindBy("height", height)
}

// GetAvgRecentTimes Gets average block times for recent blocks by limit
func (s *BlockSeq) GetAvgRecentTimes(limit int64) store.GetAvgRecentTimesResult {
	defer LogQueryDuration(time.Now(), "BlockSeqStore_GetAvgRecentTimes")

	var res store.GetAvgRecentTimesResult
	s.db.Raw(blockTimesForRecentBlocksQuery, limit).Scan(&res)

	return res
}

// FindMostRecent finds most recent block sequence
func (s *BlockSeq) FindMostRecent() (*model.BlockSeq, error) {
	blockSeq := &model.BlockSeq{}
	if err := findMostRecent(s.db, "time", blockSeq); err != nil {
		return nil, err
	}
	return blockSeq, nil
}

// DeleteOlderThan deletes block sequence older than given threshold
func (s *BlockSeq) DeleteOlderThan(purgeThreshold time.Time, activityPeriods []store.ActivityPeriodRow) (*int64, error) {
	tx := s.db.
		Unscoped()

	hasIntervals := false
	for _, activityPeriod := range activityPeriods {
		// Make sure that there are many intervals (ie. days) in period
		if !activityPeriod.Min.Equal(activityPeriod.Max) {
			hasIntervals = true
			// Thus, we do not add 1 day to Max because we don't want to purge sequences within last day of period
			tx = tx.Where("time >= ? AND time < ?", activityPeriod.Min, activityPeriod.Max)
		}
	}

	if hasIntervals {
		tx.Where("time < ?", purgeThreshold).
			Delete(&model.BlockSeq{})

		if tx.Error != nil {
			return nil, checkErr(tx.Error)
		}
	} else {
		logger.Info("no block sequences to purge")
	}

	return &tx.RowsAffected, nil
}

// Summarize gets the summarized version of block sequences
func (s *BlockSeq) Summarize(interval types.SummaryInterval, activityPeriods []store.ActivityPeriodRow) ([]store.BlockSeqSummary, error) {
	defer LogQueryDuration(time.Now(), "BlockSummaryStore_Summarize")

	tx := s.db.
		Table(model.BlockSeq{}.TableName()).
		Select(summarizeBlocksQuerySelect, interval).
		Order("time_bucket").
		Group("time_bucket")

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

	var models []store.BlockSeqSummary
	for rows.Next() {
		var summary store.BlockSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}