package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/client/figmentclient"

	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameHeightMetaRetriever = "HeightMetaRetriever"
)

var (
	_ pipeline.Task = (*heightMetaRetrieverTask)(nil)
)

func NewHeightMetaRetrieverTask(c figmentclient.Client) *heightMetaRetrieverTask {
	return &heightMetaRetrieverTask{
		client:         c,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameHeightMetaRetriever),
	}
}

type heightMetaRetrieverTask struct {
	client         figmentclient.Client
	metricObserver metrics.Observer
}

type HeightMeta struct {
	ChainId     uint64
	Height      int64
	Time        *types.Time
	Epoch       *int64
	LastInEpoch *bool
}

func (t *heightMetaRetrieverTask) GetName() string {
	return TaskNameHeightMetaRetriever
}

func (t *heightMetaRetrieverTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

	chainStatus, err := t.client.GetChainStatus(ctx)
	if err != nil {
		return err
	}

	meta, err := t.client.GetMetaByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		return err
	}

	payload.HeightMeta = HeightMeta{
		ChainId:     chainStatus.ChainId,
		Height:      meta.Height,
		Time:        types.NewTimeFromSeconds(meta.Time),
		Epoch:       meta.Epoch,
		LastInEpoch: meta.LastInEpoch,
	}
	return nil
}
