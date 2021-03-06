package account

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getDetailsUseCase struct {
	db     *psql.Store
	client figmentclient.Client
}

func NewGetDetailsUseCase(c figmentclient.Client, db *psql.Store) *getDetailsUseCase {
	return &getDetailsUseCase{
		client: c,
		db:     db,
	}
}

func (uc *getDetailsUseCase) Execute(ctx context.Context, address string, limit int64) (*DetailsView, error) {
	lastHeightAccountInfo, err := uc.client.GetAccountByAddressAndHeight(ctx, address, 0)
	if lastHeightAccountInfo == nil && err != nil {
		return nil, err
	}

	internalTransfersSent, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, indexer.OperationTypeInternalTransferSent, limit)
	if err != nil {
		return nil, err
	}

	validatorGroupVoteCastReceived, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, indexer.OperationTypeValidatorGroupVoteCastReceived, limit)
	if err != nil {
		return nil, err
	}

	validatorGroupVoteCastSent, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, indexer.OperationTypeValidatorGroupVoteCastSent, limit)
	if err != nil {
		return nil, err
	}

	goldLocked, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, figmentclient.OperationTypeGoldLocked, limit)
	if err != nil {
		return nil, err
	}

	goldUnlocked, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, figmentclient.OperationTypeGoldUnlocked, limit)
	if err != nil {
		return nil, err
	}

	goldWithdrawn, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, figmentclient.OperationTypeGoldWithdrawn, limit)
	if err != nil {
		return nil, err
	}

	accountSlashed, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, figmentclient.OperationTypeAccountSlashed, limit)
	if err != nil {
		return nil, err
	}

	validatorPaymentDistributed, err := uc.db.GetAccounts().AccountActivitySeq.FindLastByAddressAndKind(address, figmentclient.OperationTypeValidatorEpochPaymentDistributed, limit)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(address, lastHeightAccountInfo, internalTransfersSent, validatorGroupVoteCastReceived, validatorGroupVoteCastSent, goldLocked, goldUnlocked, goldWithdrawn, accountSlashed, validatorPaymentDistributed)
}
