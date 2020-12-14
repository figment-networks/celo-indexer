package indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	ValidatorAggCreatorTaskName      = "ValidatorAggCreator"
	ValidatorGroupAggCreatorTaskName = "ValidatorGroupAggCreator"
)

var (
	_ pipeline.Task = (*validatorAggCreatorTask)(nil)
	_ pipeline.Task = (*validatorGroupAggCreatorTask)(nil)
)

func NewValidatorAggCreatorTask(db *store.Store) *validatorAggCreatorTask {
	return &validatorAggCreatorTask{
		db: db,
	}
}

type validatorAggCreatorTask struct {
	db *store.Store
}

func (t *validatorAggCreatorTask) GetName() string {
	return ValidatorAggCreatorTaskName
}

func (t *validatorAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	rawValidators := payload.RawValidators

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	var newValidatorAggs []model.ValidatorAgg
	var updatedValidatorAggs []model.ValidatorAgg
	for _, rawValidator := range rawValidators {
		existing, err := t.db.ValidatorAgg.FindByAddress(rawValidator.Address)
		if err != nil {
			if err == store.ErrNotFound {
				// Create new

				validator := model.ValidatorAgg{
					Aggregate: &model.Aggregate{
						StartedAtHeight: payload.Syncable.Height,
						StartedAt:       *payload.Syncable.Time,
						RecentAtHeight:  payload.Syncable.Height,
						RecentAt:        *payload.Syncable.Time,
					},

					Address:                 rawValidator.Address,
					RecentName:              rawValidator.Name,
					RecentMetadataUrl:       rawValidator.MetadataUrl,
					RecentAsValidatorHeight: payload.Syncable.Height,
				}

				if rawValidator.Signed == nil {
					validator.AccumulatedUptime = 0
					validator.AccumulatedUptimeCount = 0
				} else {
					if *rawValidator.Signed {
						validator.AccumulatedUptime = 1
					} else {
						validator.AccumulatedUptime = 0
					}
					validator.AccumulatedUptimeCount = 1
				}

				newValidatorAggs = append(newValidatorAggs, validator)
			} else {
				return err
			}
		} else {
			// Update
			validator := &model.ValidatorAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       *payload.Syncable.Time,
				},

				RecentAsValidatorHeight: payload.Syncable.Height,
			}

			if rawValidator.Signed != nil {
				if *rawValidator.Signed {
					validator.AccumulatedUptime = existing.AccumulatedUptime + 1
				} else {
					validator.AccumulatedUptime = existing.AccumulatedUptime
				}
				validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
			}

			if rawValidator.Name != "" {
				validator.RecentName = rawValidator.Name
			}

			if rawValidator.MetadataUrl != "" {
				validator.RecentMetadataUrl = rawValidator.MetadataUrl
			}

			existing.Update(validator)

			updatedValidatorAggs = append(updatedValidatorAggs, *existing)
		}
	}
	payload.NewValidatorAggregates = newValidatorAggs
	payload.UpdatedValidatorAggregates = updatedValidatorAggs

	return nil
}

func NewValidatorGroupAggCreatorTask(db *store.Store) *validatorGroupAggCreatorTask {
	return &validatorGroupAggCreatorTask{
		db: db,
	}
}

type validatorGroupAggCreatorTask struct {
	db *store.Store
}

func (t *validatorGroupAggCreatorTask) GetName() string {
	return ValidatorGroupAggCreatorTaskName
}

func (t *validatorGroupAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	rawValidatorGroups := payload.RawValidatorGroups

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	var newValidatorGroupAggs []model.ValidatorGroupAgg
	var updatedValidatorGroupAggs []model.ValidatorGroupAgg
	for _, rawGroup := range rawValidatorGroups {
		existing, err := t.db.ValidatorGroupAgg.FindByAddress(rawGroup.Address)
		if err != nil {
			if err == store.ErrNotFound {
				// Create new

				group := model.ValidatorGroupAgg{
					Aggregate: &model.Aggregate{
						StartedAtHeight: payload.Syncable.Height,
						StartedAt:       *payload.Syncable.Time,
						RecentAtHeight:  payload.Syncable.Height,
						RecentAt:        *payload.Syncable.Time,
					},

					Address:           rawGroup.Address,
					RecentName:        rawGroup.Name,
					RecentMetadataUrl: rawGroup.MetadataUrl,
				}

				newValidatorGroupAggs = append(newValidatorGroupAggs, group)
			} else {
				return err
			}
		} else {
			// Update
			group := &model.ValidatorGroupAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       *payload.Syncable.Time,
				},
			}

			if rawGroup.Name != "" {
				group.RecentName = rawGroup.Name
			}

			if rawGroup.MetadataUrl != "" {
				group.RecentMetadataUrl = rawGroup.MetadataUrl
			}

			existing.Update(group)

			updatedValidatorGroupAggs = append(updatedValidatorGroupAggs, *existing)
		}
	}
	payload.NewValidatorGroupAggregates = newValidatorGroupAggs
	payload.UpdatedValidatorGroupAggregates = updatedValidatorGroupAggs

	return nil
}
