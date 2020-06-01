package indexing

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"

	"github.com/figment-networks/celo-indexer/client"
)

func NewValidatorFetcherTask(client *client.Client) pipeline.Task {
	return &ValidatorFetcherTask{
		client: client,
	}
}

type ValidatorFetcherTask struct {
	client *client.Client
}

func (t *ValidatorFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	validators, err := t.client.Validator.GetByHeight(payload.currentHeight)
	if err != nil {
		return err
	}

	fmt.Println(validators)

	return nil
}
