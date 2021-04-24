package indexer

import (
	"context"
	"fmt"

	"github.com/celo-org/kliento/contracts"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"

	"github.com/figment-networks/celo-indexer/model"
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

func NewValidatorAggCreatorTask(c figmentclient.Client, validatorAggDb store.ValidatorAgg) *validatorAggCreatorTask {
	return &validatorAggCreatorTask{
		client:         c,
		validatorAggDb: validatorAggDb,
	}
}

type validatorAggCreatorTask struct {
	client         figmentclient.Client
	validatorAggDb store.ValidatorAgg
}

func (t *validatorAggCreatorTask) GetName() string {
	return ValidatorAggCreatorTaskName
}

func (t *validatorAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	rawValidators := payload.RawValidators

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	shouldFetchIdentities := false
	report, ok := ctx.Value(CtxReport).(*model.Report)
	if ok {
		// No need to get metadata for all heights
		// Get metadata only for beginning of sync cycle and end
		if payload.CurrentHeight == report.StartHeight || payload.CurrentHeight == report.EndHeight {
			shouldFetchIdentities = true
		}
	}

	existingValidatorAggs, err := t.validatorAggDb.All()
	if err != nil {
		return err
	}

	existingValidatorAggsMap := make(map[string]*model.ValidatorAgg)
	for _, validatorAgg := range existingValidatorAggs {
		va := validatorAgg
		existingValidatorAggsMap[va.Address] = &va
	}

	var newValidatorAggs []model.ValidatorAgg
	var updatedValidatorAggs []model.ValidatorAgg
	for _, rawValidator := range rawValidators {
		existing, ok := existingValidatorAggsMap[rawValidator.Address]
		if !ok {
			// Create new
			validator := model.ValidatorAgg{
				Aggregate: &model.Aggregate{
					StartedAtHeight: payload.Syncable.Height,
					StartedAt:       *payload.Syncable.Time,
					RecentAtHeight:  payload.Syncable.Height,
					RecentAt:        *payload.Syncable.Time,
				},

				Address:                 rawValidator.Address,
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

			// Always get identity for new records
			identity, err := t.client.GetIdentityByHeight(ctx, validator.Address, payload.CurrentHeight)
			if err != nil {
				return err
			}
			validator.RecentName = identity.Name
			validator.RecentMetadataUrl = identity.MetadataUrl

			newValidatorAggs = append(newValidatorAggs, validator)
		} else {
			// Update
			validator := &model.ValidatorAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       *payload.Syncable.Time,
				},

				RecentAsValidatorHeight: payload.Syncable.Height,
			}

			if rawValidator.Signed == nil {
				validator.AccumulatedUptime = existing.AccumulatedUptime
				validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount
			} else {
				if *rawValidator.Signed {
					validator.AccumulatedUptime = existing.AccumulatedUptime + 1
				} else {
					validator.AccumulatedUptime = existing.AccumulatedUptime
				}
				validator.AccumulatedUptimeCount = existing.AccumulatedUptimeCount + 1
			}

			if shouldFetchIdentities {
				identity, err := t.client.GetIdentityByHeight(ctx, validator.Address, payload.CurrentHeight)
				if err != nil {
					return err
				}

				validator.RecentName = identity.Name
				validator.RecentMetadataUrl = identity.MetadataUrl
			}

			existing.Update(validator)

			updatedValidatorAggs = append(updatedValidatorAggs, *existing)
		}
	}
	payload.NewValidatorAggregates = newValidatorAggs
	payload.UpdatedValidatorAggregates = updatedValidatorAggs

	return nil
}

func NewValidatorGroupAggCreatorTask(c figmentclient.Client, validatorGroupAggDb store.ValidatorGroupAgg) *validatorGroupAggCreatorTask {
	return &validatorGroupAggCreatorTask{
		client:              c,
		validatorGroupAggDb: validatorGroupAggDb,
	}
}

type validatorGroupAggCreatorTask struct {
	client              figmentclient.Client
	validatorGroupAggDb store.ValidatorGroupAgg
}

func (t *validatorGroupAggCreatorTask) GetName() string {
	return ValidatorGroupAggCreatorTaskName
}

func (t *validatorGroupAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	rawValidatorGroups := payload.RawValidatorGroups

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	shouldFetchIdentities := false
	report, ok := ctx.Value(CtxReport).(*model.Report)
	if ok {
		// No need to get metadata for all heights
		// Get metadata only for beginning of sync cycle and end
		if payload.CurrentHeight == report.StartHeight || payload.CurrentHeight == report.EndHeight {
			shouldFetchIdentities = true
		}
	}

	existingValidatorGroupAggs, err := t.validatorGroupAggDb.All()
	if err != nil {
		return err
	}

	existingValidatorGroupAggsMap := make(map[string]*model.ValidatorGroupAgg)
	for _, validatorGroupAgg := range existingValidatorGroupAggs {
		vga := validatorGroupAgg
		existingValidatorGroupAggsMap[vga.Address] = &vga
	}

	var newValidatorGroupAggs []model.ValidatorGroupAgg
	var updatedValidatorGroupAggs []model.ValidatorGroupAgg
	for _, rawGroup := range rawValidatorGroups {
		existing, ok := existingValidatorGroupAggsMap[rawGroup.Address]
		if !ok {
			// Create new
			group := model.ValidatorGroupAgg{
				Aggregate: &model.Aggregate{
					StartedAtHeight: payload.Syncable.Height,
					StartedAt:       *payload.Syncable.Time,
					RecentAtHeight:  payload.Syncable.Height,
					RecentAt:        *payload.Syncable.Time,
				},

				Address: rawGroup.Address,
			}

			// Always get identity for new records
			identity, err := t.client.GetIdentityByHeight(ctx, group.Address, payload.CurrentHeight)
			if err != nil {
				return err
			}
			group.RecentName = identity.Name
			group.RecentMetadataUrl = identity.MetadataUrl

			newValidatorGroupAggs = append(newValidatorGroupAggs, group)
		} else {
			// Update
			group := &model.ValidatorGroupAgg{
				Aggregate: &model.Aggregate{
					RecentAtHeight: payload.Syncable.Height,
					RecentAt:       *payload.Syncable.Time,
				},
			}

			if shouldFetchIdentities {
				identity, err := t.client.GetIdentityByHeight(ctx, group.Address, payload.CurrentHeight)
				if err != nil {
					return err
				}
				group.RecentName = identity.Name
				group.RecentMetadataUrl = identity.MetadataUrl
			}

			existing.Update(group)

			updatedValidatorGroupAggs = append(updatedValidatorGroupAggs, *existing)
		}
	}
	payload.NewValidatorGroupAggregates = newValidatorGroupAggs
	payload.UpdatedValidatorGroupAggregates = updatedValidatorGroupAggs

	return nil
}

func NewProposalAggCreatorTask(proposalAggDb store.ProposalAgg) *proposalAggCreatorTask {
	return &proposalAggCreatorTask{
		proposalAggDb: proposalAggDb,
	}
}

type proposalAggCreatorTask struct {
	proposalAggDb store.ProposalAgg
}

func (t *proposalAggCreatorTask) GetName() string {
	return ProposalAggCreatorTaskName
}

func (t *proposalAggCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	parsedGovernanceLogs := payload.ParsedGovernanceLogs

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageAggregator, t.GetName(), payload.CurrentHeight))

	existingProposalAggs, _, err := t.proposalAggDb.All(0, nil)
	if err != nil {
		return err
	}

	existingProposalAggsMap := make(map[uint64]*model.ProposalAgg)
	for _, proposalAgg := range existingProposalAggs {
		pa := proposalAgg
		existingProposalAggsMap[pa.ProposalId] = &pa
	}

	var newProposalAggs []model.ProposalAgg
	var updatedProposalAggs []model.ProposalAgg
	for _, log := range parsedGovernanceLogs {
		existing, ok := existingProposalAggsMap[log.ProposalId]
		if !ok {
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
				proposal.Deposit = event.Deposit.String()
				proposal.TransactionCount = event.TransactionCount.Int64()
				proposal.RecentStage = model.ProposalStageProposed
			}

			newProposalAggs = append(newProposalAggs, proposal)
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
				if vote == model.VoteAbstain {
					toUpdate.AbstainVotesTotal += 1
					existingValue := types.NewQuantityFromString(toUpdate.AbstainVotesWeightTotal)
					if err := existingValue.Add(types.NewQuantity(event.Weight)); err != nil {
						return err
					}
					toUpdate.AbstainVotesWeightTotal = existingValue.String()
				} else if vote == model.VoteNo {
					toUpdate.NoVotesTotal += 1
					existingValue := types.NewQuantityFromString(toUpdate.NoVotesWeightTotal)
					if err := existingValue.Add(types.NewQuantity(event.Weight)); err != nil {
						return err
					}
					toUpdate.NoVotesWeightTotal = existingValue.String()
				} else if vote == model.VoteYes {
					toUpdate.YesVotesTotal += 1
					existingValue := types.NewQuantityFromString(toUpdate.YesVotesWeightTotal)
					if err := existingValue.Add(types.NewQuantity(event.Weight)); err != nil {
						return err
					}
					toUpdate.YesVotesWeightTotal = existingValue.String()
				}

				toUpdate.VotesTotal += 1
				existingValue := types.NewQuantityFromString(toUpdate.VotesWeightTotal)
				if err := existingValue.Add(types.NewQuantity(event.Weight)); err != nil {
					return err
				}
				toUpdate.VotesWeightTotal = existingValue.String()

			case figmentclient.OperationTypeProposalUpvoted:
				event := log.Details.(*contracts.GovernanceProposalUpvoted)

				existingValue := types.NewQuantityFromString(toUpdate.VotesWeightTotal)
				if err := existingValue.Add(types.NewQuantity(event.Upvotes)); err != nil {
					return err
				}
				toUpdate.VotesWeightTotal = existingValue.String()

			case figmentclient.OperationTypeProposalApproved:
				toUpdate.ApprovedAtHeight = payload.Syncable.Height
				toUpdate.ApprovedAt = *payload.Syncable.Time
				toUpdate.RecentStage = model.ProposalStageApproved

			case figmentclient.OperationTypeProposalExecuted:
				toUpdate.ExecutedAtHeight = payload.Syncable.Height
				toUpdate.ExecutedAt = *payload.Syncable.Time
				toUpdate.RecentStage = model.ProposalStageExecuted

			case figmentclient.OperationTypeProposalDequeued:
				toUpdate.DequeuedAtHeight = payload.Syncable.Height
				toUpdate.DequeuedAt = *payload.Syncable.Time
				toUpdate.RecentStage = model.ProposalStageDequeued

			case figmentclient.OperationTypeProposalExpired:
				toUpdate.ExpiredAtHeight = payload.Syncable.Height
				toUpdate.ExpiredAt = *payload.Syncable.Time
				toUpdate.RecentStage = model.ProposalStageExpired
			}

			existing.Update(toUpdate)

			updatedProposalAggs = append(updatedProposalAggs, *existing)
		}
	}

	payload.NewProposalAggregates = newProposalAggs
	payload.UpdatedProposalAggregates = updatedProposalAggs

	return nil
}
