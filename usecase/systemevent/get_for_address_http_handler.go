package systemevent

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	_ types.HttpHandler = (*getForAddressHttpHandler)(nil)
)

type getForAddressHttpHandler struct {
	db     *psql.Store
	client figmentclient.ClientIface

	useCase *getForAddressUseCase
}

func NewGetForAddressHttpHandler(db *psql.Store, c figmentclient.ClientIface) *getForAddressHttpHandler {
	return &getForAddressHttpHandler{
		db:     db,
		client: c,
	}
}

type GetForAddressRequest struct {
	Address string                 `uri:"address" binding:"required"`
	After   *int64                 `form:"after" binding:"-"`
	Kind    *model.SystemEventKind `form:"kind" binding:"-"`
}

func (h *getForAddressHttpHandler) Handle(c *gin.Context) {
	var req GetForAddressRequest
	if err := c.ShouldBindUri(&req); err != nil {
		http.BadRequest(c, errors.New("invalid address"))
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid kind or/and after"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, req.After, req.Kind)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getForAddressHttpHandler) getUseCase() *getForAddressUseCase {
	if h.useCase == nil {
		h.useCase = NewGetForAddressUseCase(h.db)
	}
	return h.useCase
}
