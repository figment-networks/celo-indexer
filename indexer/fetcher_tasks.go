package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameBlockFetcher           = "BlockFetcher"
	TaskNameValidatorsFetcher      = "ValidatorsFetcher"
	TaskNameValidatorGroupsFetcher = "ValidatorGroupsFetcher"
)

func NewBlockFetcherTask(client figmentclient.Client) pipeline.Task {
	return &BlockFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameBlockFetcher),
	}
}

type BlockFetcherTask struct {
	client         figmentclient.Client
	metricObserver metrics.Observer
}

func (t *BlockFetcherTask) GetName() string {
	return TaskNameBlockFetcher
}

func (t *BlockFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	block, err := t.client.GetBlockByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(block,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawBlock = block
	return nil
}

func NewValidatorFetcherTask(client figmentclient.Client) pipeline.Task {
	return &ValidatorsFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameValidatorsFetcher),
	}
}

type ValidatorsFetcherTask struct {
	client         figmentclient.Client
	metricObserver metrics.Observer
}

func (t *ValidatorsFetcherTask) GetName() string {
	return TaskNameValidatorsFetcher
}

func (t *ValidatorsFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	validators, err := t.client.GetValidatorsByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(validators,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "validators"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawValidators = validators
	return nil
}

func NewValidatorGroupFetcherTask(client figmentclient.Client) pipeline.Task {
	return &ValidatorGroupsFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels(TaskNameValidatorGroupsFetcher),
	}
}

type ValidatorGroupsFetcherTask struct {
	client         figmentclient.Client
	metricObserver metrics.Observer
}

func (t *ValidatorGroupsFetcherTask) GetName() string {
	return TaskNameValidatorGroupsFetcher
}

func (t *ValidatorGroupsFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	validators, err := t.client.GetValidatorGroupsByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(validators,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "validator_groups"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawValidatorGroups = validators
	return nil
}