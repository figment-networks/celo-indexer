package model

import (
	"github.com/figment-networks/celo-indexer/types"
)

type ValidatorGroupSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	Address          string         `json:"address"`
	Name             string         `json:"name"`
	MetadataUrl      string         `json:"metadata_url"`
	Commission       types.Quantity `json:"commission"`
	ActiveVotes      types.Quantity `json:"active_votes"`
	ActiveVoteUnits  types.Quantity `json:"active_vote_units"`
	PendingVotes     types.Quantity `json:"pending_votes"`
	VotingCap     types.Quantity `json:"voting_cap"`
	MembersCount     int            `json:"members_count"`
	MembersAvgSigned float64        `json:"members_avg_signed"`
}

func (ValidatorGroupSeq) TableName() string {
	return "validator_group_sequences"
}

func (s *ValidatorGroupSeq) Valid() bool {
	return s.Sequence.Valid() &&
		s.Address != ""
}

func (s *ValidatorGroupSeq) Equal(m ValidatorGroupSeq) bool {
	return s.Address == m.Address
}

func (b *ValidatorGroupSeq) Update(m ValidatorGroupSeq) {
	b.Name = m.Name
	b.MetadataUrl = m.MetadataUrl
	b.Commission = m.Commission
	b.ActiveVotes = m.ActiveVotes
	b.ActiveVotes = m.ActiveVotes
	b.ActiveVoteUnits = m.ActiveVoteUnits
	b.PendingVotes = m.PendingVotes
	b.MembersCount = m.MembersCount
	b.MembersAvgSigned = m.MembersAvgSigned
}
