package validator

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
	_ types.HttpHandler = (*getByAddressHttpHandler)(nil)
)

type getByAddressHttpHandler struct {
	db     *psql.Store
	client figmentclient.ClientIface

	useCase *getByAddressUseCase
}

func NewGetByAddressHttpHandler(db *psql.Store, c figmentclient.ClientIface) *getByAddressHttpHandler {
	return &getByAddressHttpHandler{
		db:     db,
		client: c,
	}
}

type GetByEntityUidRequest struct {
	Address        string `uri:"address" binding:"required"`
	SequencesLimit int64  `form:"sequences_limit" binding:"-"`
}

func (h *getByAddressHttpHandler) Handle(c *gin.Context) {
	var req GetByEntityUidRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid address"))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		http.BadRequest(c, errors.New("invalid sequences limit"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, req.SequencesLimit)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, resp)
}

func (h *getByAddressHttpHandler) getUseCase() *getByAddressUseCase {
	if h.useCase == nil {
		return NewGetByAddressUseCase(h.db)
	}
	return h.useCase
}
