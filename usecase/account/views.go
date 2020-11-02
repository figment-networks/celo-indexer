package account

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
)

type HeightDetailsView struct {
	GoldBalance string `json:"gold_balance"`

	TotalLockedGold          string `json:"total_locked_gold"`
	TotalNonvotingLockedGold string `json:"total_nonvoting_locked_gold"`
	StableTokenBalance       string `json:"stable_token_balance"`
}

func ToHeightDetailsView(rawAccountDetails *figmentclient.AccountDetails) *HeightDetailsView {
	view := &HeightDetailsView{
		GoldBalance: rawAccountDetails.GoldBalance.String(),
	}

	if rawAccountDetails.TotalLockedGold != nil {
		view.TotalLockedGold = rawAccountDetails.TotalLockedGold.String()
	}

	if rawAccountDetails.TotalNonvotingLockedGold != nil {
		view.TotalNonvotingLockedGold = rawAccountDetails.TotalNonvotingLockedGold.String()
	}

	if rawAccountDetails.StableTokenBalance != nil {
		view.StableTokenBalance = rawAccountDetails.StableTokenBalance.String()
	}

	return view
}
