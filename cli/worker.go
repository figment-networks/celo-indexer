package cli

import (
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/usecase"
	"github.com/figment-networks/celo-indexer/worker"
)

func startWorker(cfg *config.Config) error {
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
	dl, err := initDataLake(cfg)
	if err != nil {
		return err
	}

	workerHandlers := usecase.NewWorkerHandlers(cfg, db, client, theCeloClient, dl)

	w, err := worker.New(cfg, workerHandlers)
	if err != nil {
		return err
	}

	return w.Start()
}
