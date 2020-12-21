package model

import "github.com/figment-networks/celo-indexer/types"

type ValidatorGroupSummary struct {
	*ModelWithTimestamps
	*Summary

	Address string `json:"address"`

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

func (ValidatorGroupSummary) TableName() string {
	return "validator_group_summary"
}
