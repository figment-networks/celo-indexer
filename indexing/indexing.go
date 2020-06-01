package indexing

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"

	"github.com/figment-networks/celo-indexer/client"
)

func StartPipeline(client *client.Client) error {
	p := pipeline.New(NewPayloadFactory())

	p.SetFetcherStage(pipeline.SyncRunner(NewValidatorFetcherTask(client)))

	err := p.Start(context.Background(), NewSource(), NewSink())
	if err != nil {
		return err
	}

	return nil
}
