package systemevent

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/http"
)

var (
	_ types.HttpHandler = (*getAllHttpHandler)(nil)
)

type getAllHttpHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *getAllUseCase
}

func NewGetAllHttpHandler(db *psql.Store, c figmentclient.Client) *getAllHttpHandler {
	return &getAllHttpHandler{
		db:     db,
		client: c,
	}
}

type GetAllRequest struct {
	Page  int64 `form:"page" binding:"required"`
	Limit int64 `form:"limit" binding:"required"`
}

func (h *getAllHttpHandler) Handle(c *gin.Context) {
	var req GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid page or limit"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Page, req.Limit)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getAllHttpHandler) getUseCase() *getAllUseCase {
	if h.useCase == nil {
		h.useCase = NewGetAllUseCase(h.db)
	}
	return h.useCase
}
