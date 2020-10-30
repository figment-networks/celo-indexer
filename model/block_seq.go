package model

import "github.com/figment-networks/celo-indexer/types"

type BlockSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	// Indexed data
	TxCount         int     `json:"tx_count"`
	Size            float64 `json:"size"`
	GasUsed         uint64  `json:"gas_used"`
	TotalDifficulty uint64  `json:"total_difficulty"`
}

func (BlockSeq) TableName() string {
	return "block_sequences"
}

func (b *BlockSeq) Valid() bool {
	return b.Sequence.Valid()
}

func (b *BlockSeq) Equal(m BlockSeq) bool {
	return b.Sequence.Equal(*m.Sequence)
}

func (b *BlockSeq) Update(m BlockSeq) {
	b.TxCount = m.TxCount
	b.Size = m.Size
	b.GasUsed = m.GasUsed
	b.TotalDifficulty = m.TotalDifficulty
}
