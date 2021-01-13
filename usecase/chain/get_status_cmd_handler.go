package chain

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type GetStatusCmdHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *getStatusUseCase
}

func NewGetStatusCmdHandler(db *psql.Store, c figmentclient.Client) *GetStatusCmdHandler {
	return &GetStatusCmdHandler{
		db:     db,
		client: c,
	}
}

func (h *GetStatusCmdHandler) Handle(ctx context.Context) {
	logger.Info("chain get status use case [handler=cmd]")

	details, err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}

	fmt.Println("=== App ===")
	fmt.Println("Name:", details.AppName)
	fmt.Println("Version", details.AppVersion)
	fmt.Println("Go version", details.GoVersion)
	fmt.Println("")

	fmt.Println("=== Chain ===")
	fmt.Println("Name:", details.ChainId)
	fmt.Println("")

	if details.IndexingStarted {
		fmt.Println("=== Indexing ===")
		fmt.Println("Last index version:", details.LastIndexVersion)
		fmt.Println("Last indexed height:", details.LastIndexedHeight)
		fmt.Println("Last indexed time:", details.LastIndexedTime)
		fmt.Println("Last indexed at:", details.LastIndexedAt)
		fmt.Println("Lag behind head:", details.Lag)
	}
	fmt.Println("")
}

func (h *GetStatusCmdHandler) getUseCase() *getStatusUseCase {
	if h.useCase == nil {
		return NewGetStatusUseCase(h.db, h.client)
	}
	return h.useCase
}
