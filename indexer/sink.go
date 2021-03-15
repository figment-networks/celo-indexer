package indexer

import (
	"context"
	"fmt"

	m "github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/pkg/errors"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/metrics"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

var (
	_ pipeline.Sink = (*sink)(nil)
)

func NewSink(syncableDb store.Syncables, databaseDb store.Database, c figmentclient.Client, versionNumber int64) *sink {
	return &sink{
		syncableDb:    syncableDb,
		databaseDb:    databaseDb,
		client:        c,
		versionNumber: versionNumber,

		databaseSizeMetric: metrics.PipelineDatabaseSizeAfterHeight.WithLabels(),
		requestCountMetric: metrics.PipelineRequestCountAfterHeight.WithLabels(),
	}
}

type sink struct {
	syncableDb    store.Syncables
	databaseDb    store.Database
	client        figmentclient.Client
	versionNumber int64

	databaseSizeMetric *m.GroupGauge
	requestCountMetric *m.GroupGauge

	successCount int64
}

func (s *sink) Consume(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.DebugJSON(payload,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "sink"),
		logger.Field("height", payload.CurrentHeight),
	)

	if err := s.setProcessed(payload); err != nil {
		return err

	}

	if err := s.addMetrics(payload); err != nil {
		return err
	}

	s.successCount++

	logger.Info(fmt.Sprintf("processing completed [status=success] [height=%d]", payload.CurrentHeight))

	return nil
}

func (s *sink) setProcessed(payload *payload) error {
	payload.Syncable.MarkProcessed(s.versionNumber, s.client.GetRequestCounter().GetCounter())
	if err := s.syncableDb.Save(payload.Syncable); err != nil {
		return errors.Wrap(err, "failed saving syncable in sink")
	}
	return nil
}

func (s *sink) addMetrics(payload *payload) error {
	res, err := s.databaseDb.GetTotalSize()
	if err != nil {
		return err
	}

	s.databaseSizeMetric.Set(res.Size)
	s.requestCountMetric.Set(float64(s.client.GetRequestCounter().GetCounter()))

	return nil
}
