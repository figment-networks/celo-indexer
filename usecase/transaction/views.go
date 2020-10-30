package transaction

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
)

type ListView struct {
	Items []*figmentclient.Transaction `json:"items"`
}

func ToListView(rawTransactions []*figmentclient.Transaction) *ListView {
	return &ListView{
		Items: rawTransactions,
	}
}
