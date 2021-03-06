package block

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

var (
	_ types.HttpHandler = (*getBlockTimesHttpHandler)(nil)
)

type getBlockTimesHttpHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *getBlockTimesUseCase
}

func NewGetBlockTimesHttpHandler(db *psql.Store, client figmentclient.Client) *getBlockTimesHttpHandler {
	return &getBlockTimesHttpHandler{
		db:     db,
		client: client,
	}
}

type GetBlockTimesRequest struct {
	Limit int64 `uri:"limit" binding:"required"`
}

func (h *getBlockTimesHttpHandler) Handle(c *gin.Context) {
	var req GetBlockTimesRequest
	if err := c.ShouldBindUri(&req); err != nil {
		log.Error(err)
		http.BadRequest(c, errors.New("invalid height"))
		return
	}

	resp := h.getUseCase().Execute(req.Limit)
	http.JsonOK(c, resp)
}

func (h *getBlockTimesHttpHandler) getUseCase() *getBlockTimesUseCase {
	if h.useCase == nil {
		return NewGetBlockTimesUseCase(h.db)
	}
	return h.useCase
}
