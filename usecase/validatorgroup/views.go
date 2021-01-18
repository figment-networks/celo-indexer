package validatorgroup

import (
	"github.com/figment-networks/celo-indexer/model"
)

type AggDetailsView struct {
	*model.ModelWithTimestamps
	*model.Aggregate

	Address           string `json:"address"`
	RecentName        string `json:"recent_name"`
	RecentMetadataUrl string `json:"recent_metadata_url"`

	LastSequences      []SeqListItem              `json:"last_sequences"`
	DelegationActivity []model.AccountActivitySeq `json:"delegation_activity"`
}

func ToAggDetailsView(m *model.ValidatorGroupAgg, validatorSequences []model.ValidatorGroupSeq, delegations []model.AccountActivitySeq) *AggDetailsView {
	view := &AggDetailsView{
		ModelWithTimestamps: m.ModelWithTimestamps,
		Aggregate:           m.Aggregate,

		Address:           m.Address,
		RecentName:        m.RecentName,
		RecentMetadataUrl: m.RecentMetadataUrl,
	}

	sequenceListView := ToSeqListView(validatorSequences)

	view.LastSequences = sequenceListView.Items
	view.DelegationActivity = delegations

	return view
}

type SeqListItem struct {
	*model.Sequence

	Address          string  `json:"address"`
	Name             string  `json:"name"`
	MetadataUrl      string  `json:"metadata_url"`
	Commission       string  `json:"commission"`
	ActiveVotes      string  `json:"active_votes"`
	PendingVotes     string  `json:"pending_votes"`
	MembersCount     int     `json:"members_count"`
	MembersAvgUptime float64 `json:"members_avg_uptime"`
}

type SeqListView struct {
	Items []SeqListItem `json:"items"`
}

func ToSeqListView(validatorGroupSeqs []model.ValidatorGroupSeq) SeqListView {
	var items []SeqListItem
	for _, m := range validatorGroupSeqs {
		item := SeqListItem{
			Sequence: m.Sequence,

			Address:          m.Address,
			Commission:       m.Commission.String(),
			ActiveVotes:      m.ActiveVotes.String(),
			PendingVotes:     m.PendingVotes.String(),
			MembersCount:     m.MembersCount,
			MembersAvgUptime: m.MembersAvgSigned,
		}

		items = append(items, item)
	}

	return SeqListView{
		Items: items,
	}
}
