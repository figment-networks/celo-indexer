package account

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	db     *psql.Store
	client figmentclient.Client
}

func NewGetByHeightUseCase(db *psql.Store, c figmentclient.Client) *getByHeightUseCase {
	return &getByHeightUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getByHeightUseCase) Execute(ctx context.Context, address string, height *int64) (*HeightDetailsView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.GetCore().Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}
	lastH := mostRecentSynced.Height

	// Show last synced height, if not provided
	if height == nil {
		height = &lastH
	}

	if *height > lastH {
		return nil, errors.New("height is not indexed yet")
	}

	accountInfo, err := uc.client.GetAccountByAddressAndHeight(ctx, address, *height)
	if err != nil {
		return nil, err
	}

	accountActivitySeqs, err := uc.db.GetAccounts().AccountActivitySeq.FindByHeightAndAddress(*height, address)
	if err != nil {
		return nil, err
	}

	return ToHeightDetailsView(accountInfo, accountActivitySeqs), nil
}
