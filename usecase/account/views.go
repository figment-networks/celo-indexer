package account

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
)

type HeightDetailsView struct {
	*figmentclient.AccountDetails
}

func ToHeightDetailsView(rawAccountDetails *figmentclient.AccountDetails) *HeightDetailsView {
	return &HeightDetailsView{
		AccountDetails: rawAccountDetails,
	}
}
