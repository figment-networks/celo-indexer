package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
	"time"
)

const (
	BlockSeqCreatorTaskName          = "BlockSeqCreator"
	ValidatorSeqCreatorTaskName      = "ValidatorSeqCreator"
	ValidatorGroupSeqCreatorTaskName = "ValidatorGroupSeqCreator"
)

var (
	_ pipeline.Task = (*blockSeqCreatorTask)(nil)
	_ pipeline.Task = (*validatorSeqCreatorTask)(nil)
	_ pipeline.Task = (*validatorGroupSeqCreatorTask)(nil)
)

// NewBlockSeqCreatorTask creates block sequences
func NewBlockSeqCreatorTask(db *store.Store) *blockSeqCreatorTask {
	return &blockSeqCreatorTask{
		db: db,
	}
}

type blockSeqCreatorTask struct {
	db *store.Store
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

	blockSeq, err := t.db.BlockSeq.FindByHeight(payload.CurrentHeight)
	if err != nil {
		if err == store.ErrNotFound {
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
func NewValidatorSeqCreatorTask(cfg *config.Config, db *store.Store) *validatorSeqCreatorTask {
	return &validatorSeqCreatorTask{
		cfg: cfg,
		db:  db,
	}
}

type validatorSeqCreatorTask struct {
	cfg *config.Config
	db  *store.Store
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

	var newValidatorSeqs []model.ValidatorSeq
	var updatedValidatorSeqs []model.ValidatorSeq
	for _, rawValidatorSeq := range mappedValidatorSeqs {
		validatorSessionSeq, err := t.db.ValidatorSeq.FindByHeightAndAddress(payload.Syncable.Height, rawValidatorSeq.Address)
		if err != nil {
			if err == store.ErrNotFound {
				newValidatorSeqs = append(newValidatorSeqs, rawValidatorSeq)
				continue
			} else {
				return err
			}
		}

		validatorSessionSeq.Update(rawValidatorSeq)
		updatedValidatorSeqs = append(updatedValidatorSeqs, *validatorSessionSeq)
	}

	payload.NewValidatorSequences = newValidatorSeqs
	payload.UpdatedValidatorSequences = updatedValidatorSeqs

	return nil
}

// NewValidatorGroupSeqCreatorTask creates validator era sequences
func NewValidatorGroupSeqCreatorTask(cfg *config.Config, db *store.Store) *validatorGroupSeqCreatorTask {
	return &validatorGroupSeqCreatorTask{
		cfg: cfg,
		db:  db,
	}
}

type validatorGroupSeqCreatorTask struct {
	cfg *config.Config
	db  *store.Store
}

func (t *validatorGroupSeqCreatorTask) GetName() string {
	return ValidatorGroupSeqCreatorTaskName
}

func (t *validatorGroupSeqCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSequencer, t.GetName(), payload.CurrentHeight))

	mappedValidatorGroupSeqs, err := ToValidatorGroupSequence(payload.Syncable, payload.RawValidatorGroups)
	if err != nil {
		return err
	}

	var newValidatorGroupSeqs []model.ValidatorGroupSeq
	var updatedValidatorGroupSeqs []model.ValidatorGroupSeq
	for _, rawValidatorGroupSeq := range mappedValidatorGroupSeqs {
		validatorEraSeq, err := t.db.ValidatorGroupSeq.FindByHeightAndAddress(payload.Syncable.Height, rawValidatorGroupSeq.Address)
		if err != nil {
			if err == store.ErrNotFound {
				newValidatorGroupSeqs = append(newValidatorGroupSeqs, rawValidatorGroupSeq)
				continue
			} else {
				return err
			}
		}

		validatorEraSeq.Update(rawValidatorGroupSeq)
		updatedValidatorGroupSeqs = append(updatedValidatorGroupSeqs, *validatorEraSeq)
	}

	payload.NewValidatorGroupSequences = newValidatorGroupSeqs
	payload.UpdatedValidatorGroupSequences = updatedValidatorGroupSeqs

	return nil
}
