package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"time"
)

type ValidatorAgg interface {
	Create(*model.ValidatorAgg) error
	Save(*model.ValidatorAgg) error
	FindBy(key string, value interface{}) (*model.ValidatorAgg, error)
	FindByID(id int64) (*model.ValidatorAgg, error)
	FindByAddress(key string) (*model.ValidatorAgg, error)
	GetAllForHeightGreaterThan(height int64) ([]model.ValidatorAgg, error)
	All() ([]model.ValidatorAgg, error)
}

type ValidatorSeq interface {
	BulkUpsert(records []model.ValidatorSeq) error
	FindByHeightAndAddress(height int64, address string) (*model.ValidatorSeq, error)
	FindByHeight(h int64) ([]model.ValidatorSeq, error)
	FindMostRecent() (*model.ValidatorSeq, error)
	FindLastByAddress(address string, limit int64) ([]model.ValidatorSeq, error)
	DeleteOlderThan(purgeThreshold time.Time) (*int64, error)
	Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]ValidatorSeqSummary, error)
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