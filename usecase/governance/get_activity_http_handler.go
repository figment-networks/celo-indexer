package governance

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
	_ types.HttpHandler = (*getActivityHttpHandler)(nil)
)

type getActivityHttpHandler struct {
	db     *psql.Store
	client figmentclient.ClientIface

	useCase *getActivityUseCase
}

func NewGetActivityHttpHandler(db *psql.Store, c figmentclient.ClientIface) *getActivityHttpHandler {
	return &getActivityHttpHandler{
		db:     db,
		client: c,
	}
}

type GetActivityRequest struct {
	ProposalId uint64 `uri:"proposal_id" binding:"required"`
	Cursor     *int64 `form:"cursor" binding:"-"`
	PageSize   *int64 `form:"page_size" binding:"-"`
}

func (h *getActivityHttpHandler) Handle(c *gin.Context) {
	var req GetActivityRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid proposal id"))
		return
	}
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid cursor or page_size params"))
		return
	}

	ds, err := h.getUseCase().Execute(c, req.ProposalId, req.Cursor, req.PageSize)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, ds)
}

func (h *getActivityHttpHandler) getUseCase() *getActivityUseCase {
	if h.useCase == nil {
		return NewGetActivityUseCase(h.client, h.db)
	}
	return h.useCase
}
