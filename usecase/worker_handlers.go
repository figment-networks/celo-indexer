package usecase

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/indexing"
)

func NewWorkerHandlers(cfg *config.Config, db *store.Store, c figmentclient.Client) *WorkerHandlers {
	return &WorkerHandlers{
		RunIndexer:       indexing.NewRunWorkerHandler(cfg, db, c),
		SummarizeIndexer: indexing.NewSummarizeWorkerHandler(cfg, db, c),
		PurgeIndexer:     indexing.NewPurgeWorkerHandler(cfg, db, c),
	}
}

type WorkerHandlers struct {
	RunIndexer       types.WorkerHandler
	SummarizeIndexer types.WorkerHandler
	PurgeIndexer     types.WorkerHandler
}
