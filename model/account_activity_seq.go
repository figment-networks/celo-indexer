package model

import "github.com/figment-networks/celo-indexer/types"

type AccountActivitySeq struct {
	*Sequence

	TransactionHash string         `json:"transaction_hash"`
	Address         string         `json:"address"`
	Amount          types.Quantity `json:"amount"`
	Kind            string         `json:"kind"`
	Data            types.Jsonb    `json:"data"`
}

func (AccountActivitySeq) TableName() string {
	return "account_activity_sequences"
}

func (b *AccountActivitySeq) Update(m AccountActivitySeq) {
	b.TransactionHash = m.TransactionHash
	b.Address = m.Address
	b.Amount = m.Amount
	b.Kind = m.Kind
	b.Data = m.Data
}
