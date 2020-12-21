package cli

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/usecase"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/pkg/errors"
)

func runCmd(cfg *config.Config, flags Flags) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	client, err := initClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()
	theCeloClient, err := initTheCeloClient(cfg)
	if err != nil {
		return err
	}

	cmdHandlers := usecase.NewCmdHandlers(cfg, db, client, theCeloClient)

	logger.Info(fmt.Sprintf("executing cmd %s ...", flags.runCommand), logger.Field("app", "cli"))

	ctx := context.Background()
	switch flags.runCommand {
	case "status":
		cmdHandlers.GetStatus.Handle(ctx)
	case "indexer_start":
		cmdHandlers.StartIndexer.Handle(ctx, flags.batchSize)
	case "indexer_backfill":
		cmdHandlers.BackfillIndexer.Handle(ctx, flags.parallel, flags.force, flags.targetIds)
	case "indexer_summarize":
		cmdHandlers.SummarizeIndexer.Handle(ctx)
	case "indexer_purge":
		cmdHandlers.PurgeIndexer.Handle(ctx)
	case "update_proposals":
		cmdHandlers.UpdateProposals.Handle(ctx)
	default:
		return errors.New(fmt.Sprintf("command %s not found", flags.runCommand))
	}
	return nil
}

