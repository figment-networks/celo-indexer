package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameBlockFetcher           = "BlockFetcher"
	TaskNameValidatorsFetcher      = "ValidatorsFetcher"
	TaskNameValidatorGroupsFetcher = "ValidatorGroupsFetcher"
	TaskNameTransactionsFetcher    = "TransactionsFetcher"
)

func NewBlockFetcherTask(client figmentclient.Client) pipeline.Task {
	return &BlockFetcherTask{client: client}
}

type BlockFetcherTask struct {
	client figmentclient.Client
}

func (t *BlockFetcherTask) GetName() string {
	return TaskNameBlockFetcher
}

func (t *BlockFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	block, err := t.client.GetBlockByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		if err == figmentclient.ErrContractNotDeployed {
			logger.Info(err.Error())
		} else {
			return err
		}
	}

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
	return &ValidatorsFetcherTask{client: client}
}

type ValidatorsFetcherTask struct {
	client figmentclient.Client
}

func (t *ValidatorsFetcherTask) GetName() string {
	return TaskNameValidatorsFetcher
}

func (t *ValidatorsFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	validators, err := t.client.GetValidatorsByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		if err == figmentclient.ErrContractNotDeployed {
			logger.Info(err.Error())
		} else {
			return err
		}
	}

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
	return &ValidatorGroupsFetcherTask{client: client}
}

type ValidatorGroupsFetcherTask struct {
	client figmentclient.Client
}

func (t *ValidatorGroupsFetcherTask) GetName() string {
	return TaskNameValidatorGroupsFetcher
}

func (t *ValidatorGroupsFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	validators, err := t.client.GetValidatorGroupsByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		if err == figmentclient.ErrContractNotDeployed {
			logger.Info(err.Error())
		} else {
			return err
		}
	}

	logger.DebugJSON(validators,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "validator_groups"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawValidatorGroups = validators
	return nil
}

func NewTransactionFetcherTask(client figmentclient.Client) pipeline.Task {
	return &TransactionsFetcherTask{client: client}
}

type TransactionsFetcherTask struct {
	client figmentclient.Client
}

func (t *TransactionsFetcherTask) GetName() string {
	return TaskNameTransactionsFetcher
}

func (t *TransactionsFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	transactions, err := t.client.GetTransactionsByHeight(ctx, payload.CurrentHeight)
	if err != nil {
		if err == figmentclient.ErrContractNotDeployed {
			logger.Info(err.Error())
		} else {
			return err
		}
	}

	logger.DebugJSON(transactions,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "transactions"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawTransactions = transactions
	return nil
}
