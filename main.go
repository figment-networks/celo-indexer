package main

import (
	"github.com/figment-networks/celo-indexer/client"
	"github.com/figment-networks/celo-indexer/indexing"
)

func main() {
	client, err := client.New("localhost:50051")
	if err != nil {
		panic(err)
	}

	err = indexing.StartPipeline(client)
	if err != nil {
		panic(err)
	}
}
