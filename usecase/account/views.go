package account

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/model"
)

type HeightDetailsView struct {
	*IdentityDetails
	*BalanceDetails

	Activity []model.AccountActivitySeq `json:"activity"`
}

func ToHeightDetailsView(rawAccountInfo *figmentclient.AccountInfo, accountActivitySeqs []model.AccountActivitySeq) *HeightDetailsView {
	view := &HeightDetailsView{
		IdentityDetails: ToIdentityDetails(rawAccountInfo),
		BalanceDetails:  ToBalanceDetails(rawAccountInfo),
	}

	view.Activity = accountActivitySeqs

	return view
}

type DetailsView struct {
	Address string `json:"address"`

	Type        string `json:"type"`
	Affiliation string `json:"affiliation"`

	*IdentityDetails
	*BalanceDetails

	InternalTransfersSent          []model.AccountActivitySeq `json:"internal_transfers_sent"`
	ValidatorGroupVoteCastReceived []model.AccountActivitySeq `json:"validator_group_vote_cast_received"`
	ValidatorGroupVoteCastSent     []model.AccountActivitySeq `json:"validator_group_vote_cast_sent"`
	GoldLocked                     []model.AccountActivitySeq `json:"gold_locked"`
	GoldUnlocked                   []model.AccountActivitySeq `json:"gold_unlocked"`
	GoldWithdrawn                  []model.AccountActivitySeq `json:"gold_withdrawn"`
	AccountSlashed                 []model.AccountActivitySeq `json:"account_slashed"`
	RewardReceived                 []model.AccountActivitySeq `json:"reward_received"`
}

func ToDetailsView(address string, rawAccountInfo *figmentclient.AccountInfo, internalTransfersSent, validatorGroupVoteCastReceived, validatorGroupVoteCastSent, goldLocked, goldUnlocked, goldWithdrawn, accountSlashed, rewardReceived []model.AccountActivitySeq) (*DetailsView, error) {
	view := &DetailsView{
		Address:         address,
		IdentityDetails: ToIdentityDetails(rawAccountInfo),
		BalanceDetails:  ToBalanceDetails(rawAccountInfo),
	}

	view.InternalTransfersSent = internalTransfersSent
	view.ValidatorGroupVoteCastReceived = validatorGroupVoteCastReceived
	view.ValidatorGroupVoteCastSent = validatorGroupVoteCastSent
	view.GoldLocked = goldLocked
	view.GoldUnlocked = goldUnlocked
	view.GoldWithdrawn = goldWithdrawn
	view.AccountSlashed = accountSlashed
	view.RewardReceived = rewardReceived

	return view, nil
}

type IdentityDetails struct {
	Name        string `json:"name"`
	MetadataUrl string `json:"metadata_url"`
	Type        string `json:"type"`
	Affiliation string `json:"affiliation"`
}

type BalanceDetails struct {
	GoldBalance string `json:"gold_balance"`

	TotalLockedGold          string `json:"total_locked_gold"`
	TotalNonvotingLockedGold string `json:"total_nonvoting_locked_gold"`
	StableTokenBalance       string `json:"stable_token_balance"`
}

func ToIdentityDetails(rawAccountInfo *figmentclient.AccountInfo) *IdentityDetails {
	return &IdentityDetails{
		Name:        rawAccountInfo.Name,
		MetadataUrl: rawAccountInfo.MetadataUrl,
		Type:        rawAccountInfo.Type,
		Affiliation: rawAccountInfo.Affiliation,
	}
}

func ToBalanceDetails(rawAccountInfo *figmentclient.AccountInfo) *BalanceDetails {
	details := &BalanceDetails{}

	details.GoldBalance = rawAccountInfo.GoldBalance.String()
	if rawAccountInfo.TotalLockedGold != nil {
		details.TotalLockedGold = rawAccountInfo.TotalLockedGold.String()
	}

	if rawAccountInfo.TotalNonvotingLockedGold != nil {
		details.TotalNonvotingLockedGold = rawAccountInfo.TotalNonvotingLockedGold.String()
	}

	if rawAccountInfo.StableTokenBalance != nil {
		details.StableTokenBalance = rawAccountInfo.StableTokenBalance.String()
	}

	return details
}
