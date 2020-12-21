package governance

import "github.com/figment-networks/celo-indexer/model"

type ProposalListView struct {
	Items      []model.ProposalAgg `json:"items"`
	NextCursor *int64              `json:"next_cursor,omitempty"`
}

func ToProposalListView(proposalAggs []model.ProposalAgg, nextCursor *int64) *ProposalListView {
	view := &ProposalListView{
		Items: proposalAggs,
		NextCursor: nextCursor,
	}

	return view
}

type ActivityListView struct {
	Items      []model.GovernanceActivitySeq `json:"items"`
	NextCursor *int64                        `json:"next_cursor,omitempty"`
}

func ToActivityListView(governanceActivitySeqs []model.GovernanceActivitySeq, nextCursor *int64) *ActivityListView {

	view := &ActivityListView{
		Items:      governanceActivitySeqs,
		NextCursor: nextCursor,
	}

	return view
}
