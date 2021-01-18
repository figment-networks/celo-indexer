package model

import "github.com/figment-networks/celo-indexer/types"

type GovernanceActivitySeq struct {
	*Model
	*Sequence

	ProposalId      uint64      `json:"proposal_id"`
	Account         string      `json:"account"`
	TransactionHash string      `json:"transaction_hash"`
	Kind            string      `json:"kind"`
	Data            types.Jsonb `json:"data"`
}

func (GovernanceActivitySeq) TableName() string {
	return "governance_activity_sequences"
}

func (b *GovernanceActivitySeq) Update(m GovernanceActivitySeq) {
	b.Kind = m.Kind
	b.Data = m.Data
}
