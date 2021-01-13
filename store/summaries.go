package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"time"
)

type BlockSummary interface {
	Find(query *model.BlockSummary) (*model.BlockSummary, error)
	FindMostRecent() (*model.BlockSummary, error)
	FindMostRecentByInterval(interval types.SummaryInterval) (*model.BlockSummary, error)
	FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error)
	FindSummary(interval types.SummaryInterval, period string) ([]model.BlockSummary, error)
	DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error)
}

type ValidatorSummary interface {
	Find(query *model.ValidatorSummary) (*model.ValidatorSummary, error)
	FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error)
	FindSummary(interval types.SummaryInterval, period string) ([]ValidatorSummaryRow, error)
	FindSummaryByAddress(address string, interval types.SummaryInterval, period string) ([]model.ValidatorSummary, error)
	FindMostRecent() (*model.ValidatorSummary, error)
	FindMostRecentByInterval(interval types.SummaryInterval) (*model.ValidatorSummary, error)
	DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error)
}

type ValidatorGroupSummary interface {
	Find(query *model.ValidatorGroupSummary) (*model.ValidatorGroupSummary, error)
	FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error)
	FindSummary(interval types.SummaryInterval, period string) ([]ValidatorGroupSummaryRow, error)
	FindSummaryByAddress(address string, interval types.SummaryInterval, period string) ([]model.ValidatorGroupSummary, error)
	FindMostRecent() (*model.ValidatorGroupSummary, error)
	FindMostRecentByInterval(interval types.SummaryInterval) (*model.ValidatorGroupSummary, error)
	DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error)
}

type ActivityPeriodRow struct {
	Period int64
	Min    types.Time
	Max    types.Time
}

type ValidatorSummaryRow struct {
	TimeBucket   string  `json:"time_bucket"`
	TimeInterval string  `json:"time_interval"`
	SignedAvg    float64 `json:"signed_avg"`
	ScoreAvg     string  `json:"score_avg"`
	ScoreMin     string  `json:"score_min"`
	ScoreMax     string  `json:"score_max"`
}

type ValidatorGroupSummaryRow struct {
	TimeBucket         string `json:"time_bucket"`
	TimeInterval       string `json:"time_interval"`
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
