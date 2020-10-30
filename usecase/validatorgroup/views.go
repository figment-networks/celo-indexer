package validatorgroup

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
)

type AggDetailsView struct {
	*model.Model
	*model.Aggregate

	Address string `json:"address"`

	LastSequences []model.ValidatorGroupSeq `json:"last_sequences"`
}

func ToAggDetailsView(m *model.ValidatorGroupAgg, validatorSequences []model.ValidatorGroupSeq) *AggDetailsView {
	return &AggDetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		Address: m.Address,

		LastSequences: validatorSequences,
	}
}

type SeqListItem struct {
	*model.Sequence

	Address         string         `json:"address"`
	Commission      types.Quantity `json:"commission"`
	ActiveVotes     types.Quantity `json:"active_votes"`
	ActiveVoteUnits types.Quantity `json:"active_vote_units"`
	PendingVotes    types.Quantity `json:"pending_votes"`
	MembersCount    int            `json:"members_count"`
}

type SeqListView struct {
	Items []SeqListItem `json:"items"`
}

func ToSeqListView(validatorSeqs []model.ValidatorGroupSeq) SeqListView {
	var items []SeqListItem
	for _, m := range validatorSeqs {
		item := SeqListItem{
			Sequence: m.Sequence,

			Address:         m.Address,
			Commission:      m.Commission,
			ActiveVotes:     m.ActiveVotes,
			ActiveVoteUnits: m.ActiveVoteUnits,
			PendingVotes:    m.PendingVotes,
			MembersCount:    m.MembersCount,
		}

		items = append(items, item)
	}

	return SeqListView{
		Items: items,
	}
}
