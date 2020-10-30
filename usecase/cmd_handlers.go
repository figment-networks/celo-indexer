package usecase

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/usecase/chain"
	"github.com/figment-networks/celo-indexer/usecase/indexing"
)

func NewCmdHandlers(cfg *config.Config, db *store.Store, c figmentclient.Client) *CmdHandlers {
	return &CmdHandlers{
		GetStatus:        chain.NewGetStatusCmdHandler(db, c),
		StartIndexer:     indexing.NewStartCmdHandler(cfg, db, c),
		BackfillIndexer:  indexing.NewBackfillCmdHandler(cfg, db, c),
		PurgeIndexer:     indexing.NewPurgeCmdHandler(cfg, db, c),
		SummarizeIndexer: indexing.NewSummarizeCmdHandler(cfg, db, c),
	}
}

type CmdHandlers struct {
	GetStatus        *chain.GetStatusCmdHandler
	StartIndexer     *indexing.StartCmdHandler
	BackfillIndexer  *indexing.BackfillCmdHandler
	PurgeIndexer     *indexing.PurgeCmdHandler
	SummarizeIndexer *indexing.SummarizeCmdHandler
}
