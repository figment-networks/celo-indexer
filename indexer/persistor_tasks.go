package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
	"time"
)

const (
	SyncerPersistorTaskName            = "SyncerPersistor"
	BlockSeqPersistorTaskName          = "BlockSeqPersistor"
	ValidatorSeqPersistorTaskName      = "ValidatorSeqPersistor"
	ValidatorGroupSeqPersistorTaskName = "ValidatorGroupSeqPersistor"
	ValidatorAggPersistorTaskName      = "ValidatorAggPersistor"
	ValidatorGroupAggPersistorTaskName = "ValidatorGroupAggPersistor"
)

// NewSyncerPersistorTask is responsible for storing syncable to persistence layer
func NewSyncerPersistorTask(db *store.Store) pipeline.Task {
	return &syncerPersistorTask{
		db: db,
	}
}

type syncerPersistorTask struct {
	db *store.Store
}

func (t *syncerPersistorTask) GetName() string {
	return SyncerPersistorTaskName
}

func (t *syncerPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	return t.db.Syncables.CreateOrUpdate(payload.Syncable)
}

// NewBlockSeqPersistorTask is responsible for storing block to persistence layer
func NewBlockSeqPersistorTask(db *store.Store) pipeline.Task {
	return &blockSeqPersistorTask{
		db: db,
	}
}

type blockSeqPersistorTask struct {
	db *store.Store
}

func (t *blockSeqPersistorTask) GetName() string {
	return BlockSeqPersistorTaskName
}

func (t *blockSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	if payload.NewBlockSequence != nil {
		return t.db.BlockSeq.Create(payload.NewBlockSequence)
	}

	if payload.UpdatedBlockSequence != nil {
		return t.db.BlockSeq.Save(payload.UpdatedBlockSequence)
	}

	return nil
}

// NewValidatorSeqPersistorTask is responsible for storing validator info to persistence layer
func NewValidatorSeqPersistorTask(db *store.Store) pipeline.Task {
	return &validatorSeqPersistorTask{
		db: db,
	}
}

type validatorSeqPersistorTask struct {
	db *store.Store
}

func (t *validatorSeqPersistorTask) GetName() string {
	return ValidatorSeqPersistorTaskName
}

func (t *validatorSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, sequence := range payload.NewValidatorSequences {
		if err := t.db.ValidatorSeq.Create(&sequence); err != nil {
			return err
		}
	}

	for _, sequence := range payload.UpdatedValidatorSequences {
		if err := t.db.ValidatorSeq.Save(&sequence); err != nil {
			return err
		}
	}

	return nil
}

// NewValidatorGroupSeqPersistorTask is responsible for storing validator era info to persistence layer
func NewValidatorGroupSeqPersistorTask(db *store.Store) pipeline.Task {
	return &validatorEraSeqPersistorTask{
		db: db,
	}
}

type validatorEraSeqPersistorTask struct {
	db *store.Store
}

func (t *validatorEraSeqPersistorTask) GetName() string {
	return ValidatorGroupSeqPersistorTaskName
}

func (t *validatorEraSeqPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, sequence := range payload.NewValidatorGroupSequences {
		if err := t.db.ValidatorGroupSeq.Create(&sequence); err != nil {
			return err
		}
	}

	for _, sequence := range payload.UpdatedValidatorGroupSequences {
		if err := t.db.ValidatorGroupSeq.Save(&sequence); err != nil {
			return err
		}
	}

	return nil
}

func NewValidatorAggPersistorTask(db *store.Store) pipeline.Task {
	return &validatorAggPersistorTask{
		db: db,
	}
}

type validatorAggPersistorTask struct {
	db *store.Store
}

func (t *validatorAggPersistorTask) GetName() string {
	return ValidatorAggPersistorTaskName
}

func (t *validatorAggPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, aggregate := range payload.NewValidatorAggregates {
		if err := t.db.ValidatorAgg.Create(&aggregate); err != nil {
			return err
		}
	}

	for _, aggregate := range payload.UpdatedValidatorAggregates {
		if err := t.db.ValidatorAgg.Save(&aggregate); err != nil {
			return err
		}
	}

	return nil
}

func NewValidatorGroupAggPersistorTask(db *store.Store) pipeline.Task {
	return &validatorGroupAggPersistorTask{
		db: db,
	}
}

type validatorGroupAggPersistorTask struct {
	db *store.Store
}

func (t *validatorGroupAggPersistorTask) GetName() string {
	return ValidatorGroupAggPersistorTaskName
}

func (t *validatorGroupAggPersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StagePersistor, t.GetName(), payload.CurrentHeight))

	for _, aggregate := range payload.NewValidatorGroupAggregates {
		if err := t.db.ValidatorGroupAgg.Create(&aggregate); err != nil {
			return err
		}
	}

	for _, aggregate := range payload.UpdatedValidatorGroupAggregates {
		if err := t.db.ValidatorGroupAgg.Save(&aggregate); err != nil {
			return err
		}
	}

	return nil
}