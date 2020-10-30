package chain

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
)

type DetailsView struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	GoVersion  string `json:"go_version"`

	ChainId uint64 `json:"chain_id,omitempty"`

	IndexingStarted   bool       `json:"indexing_started"`
	LastIndexVersion  int64      `json:"last_index_version,omitempty"`
	LastIndexedHeight int64      `json:"last_indexed_height,omitempty"`
	LastIndexedTime   types.Time `json:"last_indexed_time,omitempty"`
	LastIndexedAt     types.Time `json:"last_indexed_at,omitempty"`
	Lag               int64      `json:"indexing_lag,omitempty"`
}

func ToDetailsView(recentSyncable *model.Syncable, rawChainStatus *figmentclient.ChainStatus) *DetailsView {
	view := &DetailsView{
		AppName:    config.AppName,
		AppVersion: config.AppVersion,
		GoVersion:  config.GoVersion,
		ChainId:    rawChainStatus.ChainId,
	}

	view.IndexingStarted = recentSyncable != nil
	if view.IndexingStarted {
		view.LastIndexVersion = recentSyncable.IndexVersion
		view.LastIndexedHeight = recentSyncable.Height
		view.LastIndexedTime = *recentSyncable.Time
		view.LastIndexedAt = recentSyncable.CreatedAt

		view.Lag = rawChainStatus.LastBlockHeight - recentSyncable.Height
	}

	return view
}
