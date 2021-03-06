package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"time"
)

type BlockSeq interface {
	Create(block *model.BlockSeq) error
	Save(block *model.BlockSeq) error
	CreateIfNotExists(block *model.BlockSeq) error
	FindBy(key string, value interface{}) (*model.BlockSeq, error)
	FindByID(id int64) (*model.BlockSeq, error)
	FindByHeight(height int64) (*model.BlockSeq, error)
	GetAvgRecentTimes(limit int64) GetAvgRecentTimesResult
	FindMostRecent() (*model.BlockSeq, error)
	DeleteOlderThan(purgeThreshold time.Time, activityPeriods []ActivityPeriodRow) (*int64, error)
	Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]BlockSeqSummary, error)
}

// GetAvgRecentTimesResult Contains results for GetAvgRecentTimes query
type GetAvgRecentTimesResult struct {
	StartHeight int64   `json:"start_height"`
	EndHeight   int64   `json:"end_height"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Count       int64   `json:"count"`
	Diff        float64 `json:"diff"`
	Avg         float64 `json:"avg"`
}

type BlockSeqSummary struct {
	TimeBucket                 types.Time `json:"time_bucket"`
	Count                      int64      `json:"count"`
	BlockTimeAvg               float64    `json:"block_time_avg"`
}