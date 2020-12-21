package indexer

import (
	"context"
	"fmt"
	"github.com/celo-org/kliento/contracts"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/types"
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
	ProposalAggCreatorTaskName       = "ProposalAggCreator"
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

func NewProposalAggCreatorTask(db *store.Store) *proposalAggCreatorTask {
	return &proposalAggCreatorTask{
		db: db,
	}
}

type proposalAggCreatorTask struct {
	db *store.Store
}

func (t *proposalAggCreatorTask) GetName() string {
	return ProposalAggCreatorTaskName
}

func (t *proposalAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	parsedGovernanceLogs := payload.ParsedGovernanceLogs

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	var newProposalAggs []model.ProposalAgg
	var updatedProposalAggs []model.ProposalAgg
	for _, log := range parsedGovernanceLogs {
		existing, err := t.db.ProposalAgg.FindByProposalId(log.ProposalId)
		if err != nil {
			if err == store.ErrNotFound {
				// Create new

				proposal := model.ProposalAgg{
					Aggregate: &model.Aggregate{
						StartedAtHeight: payload.Syncable.Height,
						StartedAt:       *payload.Syncable.Time,
						RecentAtHeight:  payload.Syncable.Height,
						RecentAt:        *payload.Syncable.Time,
					},

					ProposalId: log.ProposalId,
				}

				if log.Kind == figmentclient.OperationTypeProposalQueued {
					event := log.Details.(*contracts.GovernanceProposalQueued)

					proposal.ProposedAtHeight = payload.Syncable.Height
					proposal.ProposedAt = *payload.Syncable.Time
					proposal.ProposerAddress = event.Proposer.String()
					proposal.Deposit = types.NewQuantity(event.Deposit)
					proposal.TransactionCount = event.TransactionCount.Int64()
				} else {

				}

				newProposalAggs = append(newProposalAggs, proposal)
			} else {
				return err
			}
		} else {
			// Update
			toUpdate := &model.ProposalAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       *payload.Syncable.Time,
				},
			}

			switch log.Kind {

			case figmentclient.OperationTypeProposalVoted:
				event := log.Details.(*contracts.GovernanceProposalVoted)

				vote := event.Value.Uint64()

				//According to: https://github.com/celo-org/celo-monorepo/blob/master/packages/protocol/contracts/governance/Governance.sol#L41
				if vote == 1 {
					toUpdate.AbstainVotesTotal += 1
					if err := toUpdate.AbstainVotesWeightTotal.Add(types.NewQuantity(event.Weight)); err != nil {
						return err
					}
				} else if vote == 2 {
					toUpdate.NoVotesTotal += 1
					if err := toUpdate.NoVotesWeightTotal.Add(types.NewQuantity(event.Weight)); err != nil {
						return err
					}
				} else if vote == 3 {
					toUpdate.YesVotesTotal += 1
					if err := toUpdate.YesVotesWeightTotal.Add(types.NewQuantity(event.Weight)); err != nil {
						return err
					}
				}

				toUpdate.VotesTotal += 1
				if err := toUpdate.VotesWeightTotal.Add(types.NewQuantity(event.Weight)); err != nil {
					return err
				}

			case figmentclient.OperationTypeProposalUpvoted:
				event := log.Details.(*contracts.GovernanceProposalUpvoted)

				if err := toUpdate.UpvotesTotal.Add(types.NewQuantity(event.Upvotes)); err != nil {
					return err
				}

			case figmentclient.OperationTypeProposalApproved:
				toUpdate.ApprovedAtHeight = payload.Syncable.Height
				toUpdate.ApprovedAt = *payload.Syncable.Time

			case figmentclient.OperationTypeProposalExecuted:
				toUpdate.ExecutedAtHeight = payload.Syncable.Height
				toUpdate.ExecutedAt = *payload.Syncable.Time

			case figmentclient.OperationTypeProposalDequeued:
				toUpdate.DequeuedAtHeight = payload.Syncable.Height
				toUpdate.DequeuedAt = *payload.Syncable.Time

			case figmentclient.OperationTypeProposalExpired:
				toUpdate.ExpiredAtHeight = payload.Syncable.Height
				toUpdate.ExpiredAt = *payload.Syncable.Time
			}

			existing.Update(toUpdate)

			updatedProposalAggs = append(updatedProposalAggs, *existing)
		}
	}

	payload.NewProposalAggregates = newProposalAggs
	payload.UpdatedProposalAggregates = updatedProposalAggs

	return nil
}
