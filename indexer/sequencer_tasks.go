package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
	"time"
)

const (
	BlockSeqCreatorTaskName              = "BlockSeqCreator"
	ValidatorSeqCreatorTaskName          = "ValidatorSeqCreator"
	ValidatorGroupSeqCreatorTaskName     = "ValidatorGroupSeqCreator"
	AccountActivitySeqCreatorTaskName    = "AccountActivitySeqCreator"
	GovernanceActivitySeqCreatorTaskName = "GovernanceActivitySeqCreator"
)

var (
	_ pipeline.Task = (*blockSeqCreatorTask)(nil)
	_ pipeline.Task = (*validatorSeqCreatorTask)(nil)
	_ pipeline.Task = (*validatorGroupSeqCreatorTask)(nil)
	_ pipeline.Task = (*accountActivitySeqCreatorTask)(nil)
	_ pipeline.Task = (*governanceActivitySeqCreatorTask)(nil)
)

// NewBlockSeqCreatorTask creates block sequences
func NewBlockSeqCreatorTask(blockSeqDb store.BlockSeq) *blockSeqCreatorTask {
	return &blockSeqCreatorTask{
		blockSeqDb: blockSeqDb,
	}
}

type blockSeqCreatorTask struct {
	blockSeqDb store.BlockSeq
}

func (t *blockSeqCreatorTask) GetName() string {
	return BlockSeqCreatorTaskName
}

func (t *blockSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	mappedBlockSeq, err := ToBlockSequence(payload.Syncable, payload.RawBlock)
	if err != nil {
		return err
	}

	blockSeq, err := t.blockSeqDb.FindByHeight(payload.CurrentHeight)
	if err != nil {
		if err == psql.ErrNotFound {
			payload.NewBlockSequence = mappedBlockSeq
			return nil
		} else {
			return err
		}
	}

	blockSeq.Update(*mappedBlockSeq)
	payload.UpdatedBlockSequence = blockSeq

	return nil
}

// NewValidatorSeqCreatorTask creates validator sequences
func NewValidatorSeqCreatorTask(cfg *config.Config) *validatorSeqCreatorTask {
	return &validatorSeqCreatorTask{
		cfg:            cfg,
	}
}

type validatorSeqCreatorTask struct {
	cfg            *config.Config
	validatorSeqDb store.ValidatorSeq
}

func (t *validatorSeqCreatorTask) GetName() string {
	return ValidatorSeqCreatorTaskName
}

func (t *validatorSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	mappedValidatorSeqs, err := ToValidatorSequence(payload.Syncable, payload.RawValidators)
	if err != nil {
		return err
	}

	payload.ValidatorSequences = mappedValidatorSeqs

	return nil
}

// NewValidatorGroupSeqCreatorTask creates validator era sequences
func NewValidatorGroupSeqCreatorTask(cfg *config.Config) *validatorGroupSeqCreatorTask {
	return &validatorGroupSeqCreatorTask{
		cfg: cfg,
	}
}

type validatorGroupSeqCreatorTask struct {
	cfg *config.Config
}

func (t *validatorGroupSeqCreatorTask) GetName() string {
	return ValidatorGroupSeqCreatorTaskName
}

func (t *validatorGroupSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	mappedValidatorGroupSeqs, err := ToValidatorGroupSequence(payload.Syncable, payload.RawValidatorGroups, payload.RawValidators)
	if err != nil {
		return err
	}

	payload.ValidatorGroupSequences = mappedValidatorGroupSeqs

	return nil
}

// NewAccountActivitySeqCreatorTask creates account activity sequences
func NewAccountActivitySeqCreatorTask(cfg *config.Config) *accountActivitySeqCreatorTask {
	return &accountActivitySeqCreatorTask{
		cfg: cfg,
	}
}

type accountActivitySeqCreatorTask struct {
	cfg *config.Config
}

func (t *accountActivitySeqCreatorTask) GetName() string {
	return AccountActivitySeqCreatorTaskName
}

func (t *accountActivitySeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	mappedAccountActivitySeqs, err := ToAccountActivitySequence(payload.Syncable, payload.RawTransactions)
	if err != nil {
		return err
	}

	payload.AccountActivitySequences = mappedAccountActivitySeqs

	return nil
}

// NewGovernanceActivitySeqCreatorTask creates account activity sequences
func NewGovernanceActivitySeqCreatorTask(cfg *config.Config) *governanceActivitySeqCreatorTask {
	return &governanceActivitySeqCreatorTask{
		cfg: cfg,
	}
}

type governanceActivitySeqCreatorTask struct {
	cfg *config.Config
}

func (t *governanceActivitySeqCreatorTask) GetName() string {
	return GovernanceActivitySeqCreatorTaskName
}

func (t *governanceActivitySeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	mappedGovernanceActivitySeqs, err := ToGovernanceActivitySequence(payload.Syncable, payload.ParsedGovernanceLogs)
	if err != nil {
		return err
	}

	payload.GovernanceActivitySequences = mappedGovernanceActivitySeqs

	return nil
}
