package validator

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
	_ types.HttpHandler = (*getSummaryHttpHandler)(nil)

	ErrInvalidIntervalPeriod = errors.New("invalid interval and/or period")
)

type getSummaryHttpHandler struct {
	db     *store.Store
	client figmentclient.Client

	useCase *getSummaryUseCase
}

func NewGetSummaryHttpHandler(db *store.Store, c figmentclient.Client) *getSummaryHttpHandler {
	return &getSummaryHttpHandler{
		db:     db,
		client: c,
	}
}

type GetSummaryRequest struct {
	Interval types.SummaryInterval `form:"interval" binding:"required"`
	Period   string                `form:"period" binding:"required"`
	Address  string                `form:"address" binding:"-"`
}

func (h *getSummaryHttpHandler) Handle(c *gin.Context) {
	req, err := h.validateParams(c)
	if err != nil {
		logger.Error(err)
		http.BadRequest(c, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.Interval, req.Period, req.Address)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, resp)
}

func (h *getSummaryHttpHandler) validateParams(c *gin.Context) (*GetSummaryRequest, error) {
	var req GetSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if !req.Interval.Valid() {
		return nil, ErrInvalidIntervalPeriod
	}

	return &req, nil
}

func (h *getSummaryHttpHandler) getUseCase() *getSummaryUseCase {
	if h.useCase == nil {
		return NewGetSummaryUseCase(h.db)
	}
	return h.useCase
}
