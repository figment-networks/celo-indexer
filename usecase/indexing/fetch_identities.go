package indexing

import (
	"context"
	"time"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/metrics"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type fetchIdentitiesUseCase struct {
	db     *psql.Store
	client figmentclient.Client

	validatorAggDb      store.ValidatorAgg
	validatorGroupAggDb store.ValidatorGroupAgg
}

func NewFetchIdentitiesUseCase(db *psql.Store, client figmentclient.Client) *fetchIdentitiesUseCase {
	return &fetchIdentitiesUseCase{
		db:     db,
		client: client,

		validatorAggDb:      db.GetValidators().ValidatorAgg,
		validatorGroupAggDb: db.GetValidatorGroups().ValidatorGroupAgg,
	}
}

func (uc *fetchIdentitiesUseCase) Execute(ctx context.Context) error {
	defer metrics.LogUsecaseDuration(time.Now(), "fetch_identities")

	err := uc.fetchValidatorIdentities(ctx)
	if err != nil {
		return err
	}

	err = uc.fetchValidatorGroupIdentities(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (uc *fetchIdentitiesUseCase) fetchValidatorIdentities(ctx context.Context) error {
	logger.Info("fetching validator identities...")

	validators, err := uc.validatorAggDb.All()
	if err != nil {
		return err
	}

	for _, validator := range validators {
		identity, err := uc.client.GetIdentityByHeight(ctx, validator.Address, 0)
		if err != nil {
			return err
		}

		validator.RecentName = identity.Name
		validator.RecentMetadataUrl = identity.MetadataUrl

		err = uc.validatorAggDb.UpdateIdentity(&validator)
		if err != nil {
			return err
		}
	}

	logger.Info("fetched validator identities")

	return nil
}

func (uc *fetchIdentitiesUseCase) fetchValidatorGroupIdentities(ctx context.Context) error {
	logger.Info("fetching validator group identities...")

	validatorGroups, err := uc.validatorGroupAggDb.All()
	if err != nil {
		return err
	}

	for _, validatorGroup := range validatorGroups {
		identity, err := uc.client.GetIdentityByHeight(ctx, validatorGroup.Address, 0)
		if err != nil {
			return err
		}

		validatorGroup.RecentName = identity.Name
		validatorGroup.RecentMetadataUrl = identity.MetadataUrl

		err = uc.validatorGroupAggDb.UpdateIdentity(&validatorGroup)
		if err != nil {
			return err
		}
	}

	logger.Info("fetched validator group identities")

	return nil
}
