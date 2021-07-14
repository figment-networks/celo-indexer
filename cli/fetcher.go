package cli

import (
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/fetcher"
)

func startFetchWorker(cfg *config.Config) error {
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

	dl, err := initDataLake(cfg)
	if err != nil {
		return err
	}

	return fetcher.NewWorker(cfg, client, dl).Run()
}

func startFetchManager(cfg *config.Config) error {
	pool, close, err := initWorkerPool(cfg)
	if err != nil {
		return err
	}
	defer close()

	client, err := initClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	store, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()

	manager, err := fetcher.NewManager(cfg, pool, store.GetCore().Jobs, client)
	if err != nil {
		return err
	}

	return manager.Run()
}
