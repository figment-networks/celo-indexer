package model

import (
	"github.com/figment-networks/celo-indexer/types"
	"math/big"
)

type ValidatorSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	Address     string         `json:"address"`
	Affiliation string         `json:"affiliation"`
	Signed      *bool          `json:"signed"`
	Score       types.Quantity `json:"score"`

	// Join fields
	Name        string `json:"name"`
	MetadataUrl string `json:"metadata_url"`
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
	s.Affiliation = m.Affiliation
	s.Signed = m.Signed
	s.Score = m.Score
}

func (s *ValidatorSeq) ScoreAsPercentage() float64 {
	var score, _ = new(big.Float).SetString(s.Score.String())
	var divider, _ = new(big.Float).SetString("1000000000000000000000000")
	newScore := score.Quo(score, divider)

	res, _ := newScore.Float64()

	return res
}
