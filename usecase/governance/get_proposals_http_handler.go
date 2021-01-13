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
	_ types.HttpHandler = (*getProposalsHttpHandler)(nil)
)

type getProposalsHttpHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *getProposalsUseCase
}

func NewGetProposalsHttpHandler(db *psql.Store, c figmentclient.Client) *getProposalsHttpHandler {
	return &getProposalsHttpHandler{
		db:     db,
		client: c,
	}
}

type GetProposalsRequest struct {
	Cursor   *int64 `form:"cursor" binding:"-"`
	PageSize *int64 `form:"page_size" binding:"-"`
}

func (h *getProposalsHttpHandler) Handle(c *gin.Context) {
	var req GetProposalsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid address or limit"))
		return
	}

	ds, err := h.getUseCase().Execute(c, req.Cursor, req.PageSize)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, ds)
}

func (h *getProposalsHttpHandler) getUseCase() *getProposalsUseCase {
	if h.useCase == nil {
		return NewGetProposalsUseCase(h.client, h.db)
	}
	return h.useCase
}
