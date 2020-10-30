package validator

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
)

type AggListView struct {
	Items []model.ValidatorAgg `json:"items"`
}

func ToAggListView(ms []model.ValidatorAgg) *AggListView {
	return &AggListView{
		Items: ms,
	}
}

type AggDetailsView struct {
	*model.Model
	*model.Aggregate

	Address                 string  `json:"address"`
	RecentAsValidatorHeight int64   `json:"recent_as_validator_height"`
	Uptime                  float64 `json:"uptime"`

	LastSequences []model.ValidatorSeq `json:"last_sequences"`
}

func ToAggDetailsView(m *model.ValidatorAgg, validatorSequences []model.ValidatorSeq) *AggDetailsView {
	return &AggDetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		Address:                 m.Address,
		RecentAsValidatorHeight: m.RecentAsValidatorHeight,
		Uptime:                  float64(m.AccumulatedUptime) / float64(m.AccumulatedUptimeCount),

		LastSequences: validatorSequences,
	}
}

type SeqListItem struct {
	*model.Sequence

	Address     string         `json:"address"`
	Signed      *bool          `json:"signed"`
	Score       types.Quantity `json:"score"`
	Affiliation string         `json:"affiliation"`
}

type SeqListView struct {
	Items []SeqListItem `json:"items"`
}

func ToSeqListView(validatorSeqs []model.ValidatorSeq) SeqListView {
	var items []SeqListItem
	for _, m := range validatorSeqs {
		item := SeqListItem{
			Sequence: m.Sequence,

			Address:     m.Address,
			Signed:      m.Signed,
			Score:       m.Score,
			Affiliation: m.Affiliation,
		}

		items = append(items, item)
	}

	return SeqListView{
		Items: items,
	}
}
