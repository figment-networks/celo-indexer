package model

import (
	"github.com/figment-networks/celo-indexer/types"
)

type ValidatorGroupSeq struct {
	*Model
	*Sequence

	Address          string         `json:"address"`
	Commission       types.Quantity `json:"commission"`
	ActiveVotes      types.Quantity `json:"active_votes"`
	PendingVotes     types.Quantity `json:"pending_votes"`
	VotingCap        types.Quantity `json:"voting_cap"`
	MembersCount     int            `json:"members_count"`
	MembersAvgSigned float64        `json:"members_avg_signed"`

	// Join fields
	Name        string `json:"recent_name"`
	MetadataUrl string `json:"recent_metadata_url"`
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
	b.Commission = m.Commission
	b.ActiveVotes = m.ActiveVotes
	b.ActiveVotes = m.ActiveVotes
	b.PendingVotes = m.PendingVotes
	b.MembersCount = m.MembersCount
	b.MembersAvgSigned = m.MembersAvgSigned
}
