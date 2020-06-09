package indexing

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"

	"github.com/figment-networks/celo-indexer/client"
	"github.com/figment-networks/celo-indexer/store"
)

func StartPipeline(client *client.Client, store *store.Store) error {
	p := pipeline.New(NewPayloadFactory())

	p.SetFetcherStage(pipeline.SyncRunner(NewValidatorFetcherTask(client, store)))

	err := p.Start(context.Background(), NewSource(), NewSink())
	if err != nil {
		return err
	}

	return nil
}
