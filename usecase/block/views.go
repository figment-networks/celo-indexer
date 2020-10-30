package block

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
)

type DetailsView struct {
	*figmentclient.Block
}

func ToDetailsView(rawBlock *figmentclient.Block) *DetailsView {
	view := &DetailsView{
		Block: rawBlock,
	}

	return view
}
