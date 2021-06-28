package usecase

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/client/theceloclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/usecase/chain"
	"github.com/figment-networks/celo-indexer/usecase/governance"
	"github.com/figment-networks/celo-indexer/usecase/indexing"
	"github.com/figment-networks/indexing-engine/datalake"
)

func NewCmdHandlers(cfg *config.Config, db *psql.Store, nodeClient figmentclient.Client, theCeloClient theceloclient.Client, dl *datalake.DataLake) *CmdHandlers {
	return &CmdHandlers{
		GetStatus:        chain.NewGetStatusCmdHandler(db, nodeClient),
		StartIndexer:     indexing.NewStartCmdHandler(cfg, db, nodeClient, dl),
		BackfillIndexer:  indexing.NewBackfillCmdHandler(cfg, db, nodeClient),
		PurgeIndexer:     indexing.NewPurgeCmdHandler(cfg, db, nodeClient),
		SummarizeIndexer: indexing.NewSummarizeCmdHandler(cfg, db, nodeClient),
		UpdateProposals:  governance.NewUpdateProposalsCmdHandler(db, theCeloClient),
	}
}

type CmdHandlers struct {
	GetStatus        *chain.GetStatusCmdHandler
	StartIndexer     *indexing.StartCmdHandler
	BackfillIndexer  *indexing.BackfillCmdHandler
	PurgeIndexer     *indexing.PurgeCmdHandler
	SummarizeIndexer *indexing.SummarizeCmdHandler
	UpdateProposals  *governance.UpdateProposalsCmdHandler
}
