package validatorgroup

import (
	"github.com/figment-networks/celo-indexer/model"
)

type AggDetailsView struct {
	*model.Model
	*model.Aggregate

	Address           string `json:"address"`
	RecentName        string `json:"recent_name"`
	RecentMetadataUrl string `json:"recent_metadata_url"`

	LastSequences []SeqListItem `json:"last_sequences"`
}

func ToAggDetailsView(m *model.ValidatorGroupAgg, validatorSequences []model.ValidatorGroupSeq) *AggDetailsView {
	view := &AggDetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		Address:           m.Address,
		RecentName:        m.RecentName,
		RecentMetadataUrl: m.RecentMetadataUrl,
	}

	sequenceListView := ToSeqListView(validatorSequences)

	view.LastSequences = sequenceListView.Items

	return view
}

type SeqListItem struct {
	*model.Sequence

	Address         string `json:"address"`
	Name            string `json:"name"`
	MetadataUrl     string `json:"metadata_url"`
	Commission      string `json:"commission"`
	ActiveVotes     string `json:"active_votes"`
	ActiveVoteUnits string `json:"active_vote_units"`
	PendingVotes    string `json:"pending_votes"`
	MembersCount    int    `json:"members_count"`
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
			Name:            m.Name,
			MetadataUrl:     m.MetadataUrl,
			Commission:      m.Commission.String(),
			ActiveVotes:     m.ActiveVotes.String(),
			ActiveVoteUnits: m.ActiveVoteUnits.String(),
			PendingVotes:    m.PendingVotes.String(),
			MembersCount:    m.MembersCount,
		}

		items = append(items, item)
	}

	return SeqListView{
		Items: items,
	}
}
