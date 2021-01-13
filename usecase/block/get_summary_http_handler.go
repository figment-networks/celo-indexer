package block

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
	_ types.HttpHandler = (*getBlockSummaryHttpHandler)(nil)

	ErrInvalidIntervalPeriod = errors.New("invalid interval and/or period")
)

type getBlockSummaryHttpHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *getBlockSummaryUseCase
}

func NewGetBlockSummaryHttpHandler(db *psql.Store, client figmentclient.Client) *getBlockSummaryHttpHandler {
	return &getBlockSummaryHttpHandler{
		db:     db,
		client: client,
	}
}

type GetBlockTimesForIntervalRequest struct {
	Interval types.SummaryInterval `form:"interval" binding:"required"`
	Period   string                `form:"period" binding:"required"`
}

func (h *getBlockSummaryHttpHandler) Handle(c *gin.Context) {
	req, err := h.validateParams(c)
	if err != nil {
		logger.Error(err)
		http.BadRequest(c, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.Interval, req.Period)
	if err != nil {
		logger.Error(err)
		http.ServerError(c, err)
		return
	}

	http.JsonOK(c, resp)
}

func (h *getBlockSummaryHttpHandler) validateParams(c *gin.Context) (*GetBlockTimesForIntervalRequest, error) {
	var req GetBlockTimesForIntervalRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if !req.Interval.Valid() {
		return nil, ErrInvalidIntervalPeriod
	}

	return &req, nil
}

func (h *getBlockSummaryHttpHandler) getUseCase() *getBlockSummaryUseCase {
	if h.useCase == nil {
		return NewGetBlockSummaryUseCase(h.db)
	}
	return h.useCase
}
