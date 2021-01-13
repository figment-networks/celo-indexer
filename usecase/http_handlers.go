package usecase

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/account"
	"github.com/figment-networks/celo-indexer/usecase/block"
	"github.com/figment-networks/celo-indexer/usecase/chain"
	"github.com/figment-networks/celo-indexer/usecase/governance"
	"github.com/figment-networks/celo-indexer/usecase/health"
	"github.com/figment-networks/celo-indexer/usecase/systemevent"
	"github.com/figment-networks/celo-indexer/usecase/transaction"
	"github.com/figment-networks/celo-indexer/usecase/validator"
	"github.com/figment-networks/celo-indexer/usecase/validatorgroup"
)

func NewHttpHandlers(cfg *config.Config, db *psql.Store, c figmentclient.Client) *HttpHandlers {
	return &HttpHandlers{
		Health:                     health.NewHealthHttpHandler(),
		GetStatus:                  chain.NewGetStatusHttpHandler(db, c),
		GetBlockByHeight:           block.NewGetByHeightHttpHandler(db, c),
		GetBlockTimes:              block.NewGetBlockTimesHttpHandler(db, c),
		GetBlockSummary:            block.NewGetBlockSummaryHttpHandler(db, c),
		GetTransactionsByHeight:    transaction.NewGetByHeightHttpHandler(db, c),
		GetAccountByHeight:         account.NewGetByHeightHttpHandler(db, c),
		GetAccountDetails:          account.NewGetDetailsHttpHandler(db, c),
		GetValidatorsByHeight:      validator.NewGetByHeightHttpHandler(cfg, db, c),
		GetValidatorByAddress:      validator.NewGetByAddressHttpHandler(db, c),
		GetValidatorSummary:        validator.NewGetSummaryHttpHandler(db, c),
		GetValidatorsForMinHeight:  validator.NewGetForMinHeightHttpHandler(db, c),
		GetValidatorGroupsByHeight: validatorgroup.NewGetByHeightHttpHandler(cfg, db, c),
		GetValidatorGroupByAddress: validatorgroup.NewGetByAddressHttpHandler(db, c),
		GetValidatorGroupSummary:   validatorgroup.NewGetSummaryHttpHandler(db, c),
		GetSystemEventsForAddress:  systemevent.NewGetForAddressHttpHandler(db, c),
		GetProposals:               governance.NewGetProposalsHttpHandler(db, c),
		GetProposalActivity:        governance.NewGetActivityHttpHandler(db, c),
	}
}

type HttpHandlers struct {
	Health                     types.HttpHandler
	GetStatus                  types.HttpHandler
	GetBlockTimes              types.HttpHandler
	GetBlockSummary            types.HttpHandler
	GetBlockByHeight           types.HttpHandler
	GetTransactionsByHeight    types.HttpHandler
	GetAccountByHeight         types.HttpHandler
	GetAccountDetails          types.HttpHandler
	GetValidatorsByHeight      types.HttpHandler
	GetValidatorByAddress      types.HttpHandler
	GetValidatorSummary        types.HttpHandler
	GetValidatorsForMinHeight  types.HttpHandler
	GetValidatorGroupsByHeight types.HttpHandler
	GetValidatorGroupByAddress types.HttpHandler
	GetValidatorGroupSummary   types.HttpHandler
	GetSystemEventsForAddress  types.HttpHandler
	GetProposals               types.HttpHandler
	GetProposalActivity        types.HttpHandler
}
