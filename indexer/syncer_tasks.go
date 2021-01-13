package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	MainSyncerTaskName = "MainSyncer"
)

func NewMainSyncerTask(syncableDb store.Syncables) pipeline.Task {
	return &mainSyncerTask{
		syncableDb: syncableDb,
	}
}

type mainSyncerTask struct {
	syncableDb store.Syncables
}

func (t *mainSyncerTask) GetName() string {
	return MainSyncerTaskName
}

func (t *mainSyncerTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSyncer, t.GetName(), payload.CurrentHeight))

	syncable, err := t.syncableDb.FindByHeight(payload.CurrentHeight)
	if err != nil {
		if err == psql.ErrNotFound {
			syncable = &model.Syncable{
				Height: payload.CurrentHeight,
				Time:   payload.HeightMeta.Time,

				ChainId:     payload.HeightMeta.ChainId,
				Epoch:       payload.HeightMeta.Epoch,
				LastInEpoch: payload.HeightMeta.LastInEpoch,
				Status:      model.SyncableStatusRunning,
			}
		} else {
			return err
		}
	}

	syncable.StartedAt = *types.NewTimeFromTime(time.Now())

	report, ok := ctx.Value(CtxReport).(*model.Report)
	if ok {
		syncable.ReportID = report.ID
	}

	payload.Syncable = syncable
	return nil
}
