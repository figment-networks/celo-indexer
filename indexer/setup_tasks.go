package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameSetup = "Setup"
)

var (
	_ pipeline.Task = (*setupTask)(nil)
)

func NewSetupTask(c figmentclient.Client) *setupTask {
	return &setupTask{client: c}
}

type setupTask struct {
	client figmentclient.Client
}

func (t *setupTask) GetName() string {
	return TaskNameSetup
}

func (t *setupTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

	logger.Info(fmt.Sprintf("initializing requests counter [stage=%s] [task=%s] [height=%d]", pipeline.StageSetup, t.GetName(), payload.CurrentHeight))

	t.client.GetRequestCounter().InitCounter()

	return nil
}
