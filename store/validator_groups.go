package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"time"
)

type ValidatorGroupAgg interface {
	Create(*model.ValidatorGroupAgg) error
	Save(*model.ValidatorGroupAgg) error
	CreateOrUpdate(*model.ValidatorGroupAgg) error
	FindBy(key string, value interface{}) (*model.ValidatorGroupAgg, error)
	FindByID(id int64) (*model.ValidatorGroupAgg, error)
	FindByAddress(key string) (*model.ValidatorGroupAgg, error)
	All() ([]model.ValidatorGroupAgg, error)
}

type ValidatorGroupSeq interface {
	BulkUpsert(records []model.ValidatorGroupSeq) error
	FindByHeightAndAddress(height int64, address string) (*model.ValidatorGroupSeq, error)
	FindByHeight(h int64) ([]model.ValidatorGroupSeq, error)
	FindLastByAddress(address string, limit int64) ([]model.ValidatorGroupSeq, error)
	FindMostRecent() (*model.ValidatorGroupSeq, error)
	DeleteOlderThan(purgeThreshold time.Time) (*int64, error)
	Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]ValidatorGroupSeqSummary, error)
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
	PendingVotesAvg    types.Quantity `json:"pending_votes_avg"`
	PendingVotesMin    types.Quantity `json:"pending_votes_min"`
	PendingVotesMax    types.Quantity `json:"pending_votes_max"`
}