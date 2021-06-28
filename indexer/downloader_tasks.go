package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameBlockDownloader           = "BlockDownloader"
	TaskNameValidatorsDownloader      = "ValidatorsDownloader"
	TaskNameValidatorGroupsDownloader = "ValidatorGroupsDownloader"
	TaskNameTransactionsDownloader    = "TransactionsDownloader"
)

func NewBlockDownloaderTask() pipeline.Task {
	return &BlockDownloaderTask{}
}

type BlockDownloaderTask struct{}

func (t *BlockDownloaderTask) GetName() string {
	return TaskNameBlockDownloader
}

func (t *BlockDownloaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	return payload.Retrieve("block", &payload.RawBlock)
}

func NewValidatorDownloaderTask() pipeline.Task {
	return &ValidatorsDownloaderTask{}
}

type ValidatorsDownloaderTask struct{}

func (t *ValidatorsDownloaderTask) GetName() string {
	return TaskNameValidatorsDownloader
}

func (t *ValidatorsDownloaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	return payload.Retrieve("validators", &payload.RawValidators)
}

func NewValidatorGroupDownloaderTask() pipeline.Task {
	return &ValidatorGroupsDownloaderTask{}
}

type ValidatorGroupsDownloaderTask struct{}

func (t *ValidatorGroupsDownloaderTask) GetName() string {
	return TaskNameValidatorGroupsDownloader
}

func (t *ValidatorGroupsDownloaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	return payload.Retrieve("validator_groups", &payload.RawValidatorGroups)
}

func NewTransactionDownloaderTask() pipeline.Task {
	return &TransactionsDownloaderTask{}
}

type TransactionsDownloaderTask struct{}

func (t *TransactionsDownloaderTask) GetName() string {
	return TaskNameTransactionsDownloader
}

func (t *TransactionsDownloaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))

	return payload.Retrieve("transactions", &payload.RawTransactions)
}
