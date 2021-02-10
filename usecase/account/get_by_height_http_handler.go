package account

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/http"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	_ types.HttpHandler = (*getByHeightHttpHandler)(nil)
)

type getByHeightHttpHandler struct {
	db     *psql.Store
	client figmentclient.ClientIface

	useCase *getByHeightUseCase
}

func NewGetByHeightHttpHandler(db *psql.Store, c figmentclient.ClientIface) *getByHeightHttpHandler {
	return &getByHeightHttpHandler{
		db:     db,
		client: c,
	}
}

type GetByHeightRequest struct {
	Address string `uri:"address" binding:"required"`
	Height  *int64 `form:"height" binding:"-"`
}

func (h *getByHeightHttpHandler) Handle(c *gin.Context) {
	var req GetByHeightRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid address"))
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid height"))
		return
	}

	ds, err := h.getUseCase().Execute(c, req.Address, req.Height)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, ds)
}

func (h *getByHeightHttpHandler) getUseCase() *getByHeightUseCase {
	if h.useCase == nil {
		return NewGetByHeightUseCase(h.db, h.client)
	}
	return h.useCase
}
