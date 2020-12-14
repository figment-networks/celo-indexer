package model

import "github.com/figment-networks/celo-indexer/types"

type ValidatorSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	Address     string         `json:"address"`
	Name        string         `json:"name"`
	MetadataUrl string         `json:"metadata_url"`
	Affiliation string         `json:"affiliation"`
	Signed      *bool          `json:"signed"`
	Score       types.Quantity `json:"score"`
}

func (ValidatorSeq) TableName() string {
	return "validator_sequences"
}

func (s *ValidatorSeq) Valid() bool {
	return s.Sequence.Valid() &&
		s.Address != ""
}

func (s *ValidatorSeq) Equal(m ValidatorSeq) bool {
	return s.Sequence.Equal(*m.Sequence) &&
		s.Address == m.Address
}

func (s *ValidatorSeq) Update(m ValidatorSeq) {
	s.Name = m.Name
	s.MetadataUrl = m.MetadataUrl
	s.Affiliation = m.Affiliation
	s.Signed = m.Signed
	s.Score = m.Score
}
