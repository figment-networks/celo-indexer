package main

import (
	"github.com/figment-networks/celo-indexer/client"
	"github.com/figment-networks/celo-indexer/indexing"
	"github.com/figment-networks/celo-indexer/server"
	"github.com/figment-networks/celo-indexer/store"
)

func main() {
	store, err := store.New("sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	client, err := client.New("localhost:50051")
	if err != nil {
		panic(err)
	}

	err = indexing.StartPipeline(client, store)
	if err != nil {
		panic(err)
	}

	server.Run(store)
}
