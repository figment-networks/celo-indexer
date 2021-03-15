package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/client/figmentclient"

	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameHeightMetaRetriever = "HeightMetaRetriever"
)

var (
	_ pipeline.Task = (*heightMetaRetrieverTask)(nil)
)

func NewHeightMetaRetrieverTask(c figmentclient.Client) *heightMetaRetrieverTask {
	return &heightMetaRetrieverTask{client: c}
}

type heightMetaRetrieverTask struct {
	client figmentclient.Client
}

type HeightMeta struct {
	ChainId     uint64
	Height      int64
	Time        *types.Time
	Epoch       *int64
	EpochSize   *int64
	LastInEpoch *bool
}

func (t *heightMetaRetrieverTask) GetName() string {
	return TaskNameHeightMetaRetriever
}

func (t *heightMetaRetrieverTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

	logger.Info(fmt.Sprintf("initializing requests counter [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

	t.client.GetRequestCounter().InitCounter()

	heightMeta := HeightMeta{}

	chainParams, err := t.client.GetChainParams(ctx)
	if err != nil {
		if err == figmentclient.ErrContractNotDeployed {
			logger.Info("GetChainParams returned partial data")
		} else {
			return err
		}
	} else {
		// Contract dependent data
		heightMeta.EpochSize = chainParams.EpochSize
	}

	// Get chainParams partial data
	heightMeta.ChainId = chainParams.ChainId

	meta, err := t.client.GetMetaByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		if err == figmentclient.ErrContractNotDeployed {
			logger.Info("GetMetaByHeight returned partial data")
		} else {
			return err
		}
	} else {
		// Contract dependent data
		heightMeta.Epoch = meta.Epoch
		heightMeta.LastInEpoch = meta.LastInEpoch
	}

	// Get meta partial data
	heightMeta.Height = meta.Height
	heightMeta.Time = types.NewTimeFromSeconds(meta.Time)

	payload.HeightMeta = heightMeta
	return nil
}
