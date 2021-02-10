package usecase

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/client/theceloclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/governance"
	"github.com/figment-networks/celo-indexer/usecase/indexing"
)

func NewWorkerHandlers(cfg *config.Config, db *psql.Store, client figmentclient.ClientIface, theCeloClient theceloclient.Client) *WorkerHandlers {
	return &WorkerHandlers{
		RunIndexer:       indexing.NewRunWorkerHandler(cfg, db, client),
		SummarizeIndexer: indexing.NewSummarizeWorkerHandler(cfg, db, client),
		PurgeIndexer:     indexing.NewPurgeWorkerHandler(cfg, db, client),
		UpdateProposals:  governance.NewUpdateProposalsWorkerHandler(cfg, db, theCeloClient),
	}
}

type WorkerHandlers struct {
	RunIndexer       types.WorkerHandler
	SummarizeIndexer types.WorkerHandler
	PurgeIndexer     types.WorkerHandler
	UpdateProposals  types.WorkerHandler
}
