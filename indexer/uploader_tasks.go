package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameHeightMetaUploader      = "HeightMetaUploader"
	TaskNameBlockUploader           = "BlockUploader"
	TaskNameValidatorsUploader      = "ValidatorsUploader"
	TaskNameValidatorGroupsUploader = "ValidatorGroupsUploader"
	TaskNameTransactionsUploader    = "TransactionsUploader"
)

func NewHeightMetaUploaderTask() pipeline.Task {
	return &HeightMetaUploaderTask{}
}

type HeightMetaUploaderTask struct{}

func (t *HeightMetaUploaderTask) GetName() string {
	return TaskNameHeightMetaUploader
}

func (t *HeightMetaUploaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return payload.Store("height_meta", payload.HeightMeta)
}

func NewBlockUploaderTask() pipeline.Task {
	return &BlockUploaderTask{}
}

type BlockUploaderTask struct{}

func (t *BlockUploaderTask) GetName() string {
	return TaskNameBlockUploader
}

func (t *BlockUploaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return payload.Store("block", payload.RawBlock)
}

func NewValidatorUploaderTask() pipeline.Task {
	return &ValidatorsUploaderTask{}
}

type ValidatorsUploaderTask struct{}

func (t *ValidatorsUploaderTask) GetName() string {
	return TaskNameValidatorsUploader
}

func (t *ValidatorsUploaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return payload.Store("validators", payload.RawValidators)
}

func NewValidatorGroupUploaderTask() pipeline.Task {
	return &ValidatorGroupsUploaderTask{}
}

type ValidatorGroupsUploaderTask struct{}

func (t *ValidatorGroupsUploaderTask) GetName() string {
	return TaskNameValidatorGroupsUploader
}

func (t *ValidatorGroupsUploaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return payload.Store("validator_groups", payload.RawValidatorGroups)
}

func NewTransactionUploaderTask() pipeline.Task {
	return &TransactionsUploaderTask{}
}

type TransactionsUploaderTask struct{}

func (t *TransactionsUploaderTask) GetName() string {
	return TaskNameTransactionsUploader
}

func (t *TransactionsUploaderTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return payload.Store("transactions", payload.RawTransactions)
}
