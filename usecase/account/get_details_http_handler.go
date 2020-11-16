package account

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/http"
	"github.com/figment-networks/celo-indexer/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	_ types.HttpHandler = (*getDetailsHttpHandler)(nil)
)

type getDetailsHttpHandler struct {
	db     *store.Store
	client figmentclient.Client

	useCase *getDetailsUseCase
}

func NewGetDetailsHttpHandler(db *store.Store, c figmentclient.Client) *getDetailsHttpHandler {
	return &getDetailsHttpHandler{
		db: db,
		client: c,
	}
}

type GetDetailsRequest struct {
	Address string `uri:"address" binding:"required"`
	Limit   int64  `form:"limit" binding:"required"`
}

func (h *getDetailsHttpHandler) Handle(c *gin.Context) {
	var req GetDetailsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid address or limit"))
		return
	}

	ds, err := h.getUseCase().Execute(c, req.Address, req.Limit)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, ds)
}

func (h *getDetailsHttpHandler) getUseCase() *getDetailsUseCase {
	if h.useCase == nil {
		return NewGetDetailsUseCase(h.client, h.db)
	}
	return h.useCase
}
