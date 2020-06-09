package indexing

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"

	"github.com/figment-networks/celo-indexer/client"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
)

func NewValidatorFetcherTask(client *client.Client, store *store.Store) pipeline.Task {
	return &ValidatorFetcherTask{
		client: client,
		store:  store,
	}
}

type ValidatorFetcherTask struct {
	client *client.Client
	store  *store.Store
}

func (t *ValidatorFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	validators, err := t.client.Validator.GetByHeight(payload.currentHeight)
	if err != nil {
		return err
	}

	for _, v := range validators.GetValidators() {
		t.store.Db.Create(&model.Validator{
			Address: v.Address,
			Name:    v.Name,
		})
	}

	return nil
}
