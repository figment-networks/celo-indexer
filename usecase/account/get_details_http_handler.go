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
	_ types.HttpHandler = (*getDetailsHttpHandler)(nil)
)

type getDetailsHttpHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *getDetailsUseCase
}

func NewGetDetailsHttpHandler(db *psql.Store, c figmentclient.Client) *getDetailsHttpHandler {
	return &getDetailsHttpHandler{
		db:     db,
		client: c,
	}
}

type uriParams struct {
	Address string `uri:"address" binding:"required"`
}

type queryParams struct {
	Limit int64 `form:"limit" binding:"required"`
}

func (h *getDetailsHttpHandler) Handle(c *gin.Context) {
	var uri uriParams
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid address"))
		return
	}

	var params queryParams
	if err := c.ShouldBind(&params); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid limit parameter"))
		return
	}

	ds, err := h.getUseCase().Execute(c, uri.Address, params.Limit)
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
