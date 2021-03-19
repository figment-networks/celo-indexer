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

func (s *ValidatorGroupSeq) Update(m ValidatorGroupSeq) {
	s.Commission = m.Commission
	s.ActiveVotes = m.ActiveVotes
	s.ActiveVotes = m.ActiveVotes
	s.PendingVotes = m.PendingVotes
	s.MembersCount = m.MembersCount
	s.MembersAvgSigned = m.MembersAvgSigned
}

func (s *ValidatorGroupSeq) IsValidated() bool {
	return s.MembersAvgSigned > 0
}
