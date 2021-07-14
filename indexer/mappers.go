package indexer

import (
	"encoding/json"
	"fmt"

	"github.com/celo-org/kliento/contracts"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/pkg/errors"
)

var (
	OperationTypeInternalTransferReceived                     = fmt.Sprintf("%sReceived", figmentclient.OperationTypeInternalTransfer)
	OperationTypeInternalTransferSent                         = fmt.Sprintf("%sSent", figmentclient.OperationTypeInternalTransfer)
	OperationTypeValidatorGroupVoteCastReceived               = fmt.Sprintf("%sReceived", figmentclient.OperationTypeValidatorGroupVoteCast)
	OperationTypeValidatorGroupVoteCastSent                   = fmt.Sprintf("%sSent", figmentclient.OperationTypeValidatorGroupVoteCast)
	OperationTypeValidatorGroupVoteActivatedReceived          = fmt.Sprintf("%sReceived", figmentclient.OperationTypeValidatorGroupVoteActivated)
	OperationTypeValidatorGroupVoteActivatedSent              = fmt.Sprintf("%sSent", figmentclient.OperationTypeValidatorGroupVoteActivated)
	OperationTypeValidatorGroupPendingVoteRevokedReceived     = fmt.Sprintf("%sReceived", figmentclient.OperationTypeValidatorGroupPendingVoteRevoked)
	OperationTypeValidatorGroupPendingVoteRevokedSent         = fmt.Sprintf("%sSent", figmentclient.OperationTypeValidatorGroupPendingVoteRevoked)
	OperationTypeValidatorGroupActiveVoteRevokedReceived      = fmt.Sprintf("%sReceived", figmentclient.OperationTypeValidatorGroupActiveVoteRevoked)
	OperationTypeValidatorGroupActiveVoteRevokedSent          = fmt.Sprintf("%sSent", figmentclient.OperationTypeValidatorGroupActiveVoteRevoked)
	OperationTypeValidatorEpochPaymentDistributedForGroup     = fmt.Sprintf("%sForGroup", figmentclient.OperationTypeValidatorEpochPaymentDistributed)
	OperationTypeValidatorEpochPaymentDistributedForValidator = fmt.Sprintf("%sForValidator", figmentclient.OperationTypeValidatorEpochPaymentDistributed)

	ErrBlockSequenceNotValid          = errors.New("block sequence not valid")
	ErrValidatorSequenceNotValid      = errors.New("validator sequence not valid")
	ErrValidatorGroupSequenceNotValid = errors.New("validator group sequence not valid")
)

func ToBlockSequence(syncable *model.Syncable, rawBlock *figmentclient.Block) (*model.BlockSeq, error) {
	e := &model.BlockSeq{
		Sequence: &model.Sequence{
			Height: syncable.Height,
			Time:   *syncable.Time,
		},

		TxCount:         rawBlock.TxCount,
		Size:            rawBlock.Size,
		GasUsed:         rawBlock.GasUsed,
		TotalDifficulty: rawBlock.TotalDifficulty,
	}

	if !e.Valid() {
		return nil, ErrBlockSequenceNotValid
	}

	return e, nil
}

func ToValidatorSequence(syncable *model.Syncable, rawValidators []*figmentclient.Validator) ([]model.ValidatorSeq, error) {
	var validators []model.ValidatorSeq
	for _, rawValidator := range rawValidators {
		e := model.ValidatorSeq{
			Sequence: &model.Sequence{
				Height: syncable.Height,
				Time:   *syncable.Time,
			},

			Address:     rawValidator.Address,
			Affiliation: rawValidator.Affiliation,
			Signed:      rawValidator.Signed,
			Score:       types.NewQuantity(rawValidator.Score),
		}

		if !e.Valid() {
			return nil, ErrValidatorSequenceNotValid
		}

		validators = append(validators, e)
	}
	return validators, nil
}

func ToValidatorGroupSequence(syncable *model.Syncable, rawValidatorGroups []*figmentclient.ValidatorGroup, rawValidators []*figmentclient.Validator) ([]model.ValidatorGroupSeq, error) {
	membersAvgSignedMap := getMembersAvgSignedMap(rawValidators)

	var validators []model.ValidatorGroupSeq
	for _, rawValidatorGroup := range rawValidatorGroups {
		e := model.ValidatorGroupSeq{
			Sequence: &model.Sequence{
				Height: syncable.Height,
				Time:   *syncable.Time,
			},

			Address:      rawValidatorGroup.Address,
			Commission:   types.NewQuantity(rawValidatorGroup.Commission),
			ActiveVotes:  types.NewQuantity(rawValidatorGroup.ActiveVotes),
			PendingVotes: types.NewQuantity(rawValidatorGroup.PendingVotes),
			VotingCap:    types.NewQuantity(rawValidatorGroup.VotingCap),
			MembersCount: len(rawValidatorGroup.Members),
		}

		found, ok := membersAvgSignedMap[rawValidatorGroup.Address]
		if ok {
			e.MembersAvgSigned = float64(found.count) / float64(found.total)
		}

		if !e.Valid() {
			return nil, ErrValidatorGroupSequenceNotValid
		}

		validators = append(validators, e)
	}
	return validators, nil
}

type membersAvgSigned struct {
	total int64
	count int64
}

func getMembersAvgSignedMap(rawValidators []*figmentclient.Validator) map[string]membersAvgSigned {
	membersAvgSignedMap := make(map[string]membersAvgSigned)
	for _, rawValidator := range rawValidators {
		foundInfo, ok := membersAvgSignedMap[rawValidator.Affiliation]
		if ok {
			if rawValidator.Signed != nil && *rawValidator.Signed {
				foundInfo.total += 1
			}
			foundInfo.count += 1
		} else {
			membersAvgSignedMap[rawValidator.Affiliation] = membersAvgSigned{
				total: 1,
				count: 1,
			}
		}
	}
	return membersAvgSignedMap
}

func ToAccountActivitySequence(syncable *model.Syncable, rawTransactions []*figmentclient.Transaction) ([]model.AccountActivitySeq, error) {
	var accountActivities []model.AccountActivitySeq
	for _, rawTransaction := range rawTransactions {
		for _, rawOperation := range rawTransaction.Operations {
			sequence := &model.Sequence{
				Height: syncable.Height,
				Time:   *syncable.Time,
			}

			accountActivitiesFromOperation, err := operationToAccountActivitySequence(sequence, rawOperation, rawTransaction.Hash)
			if err != nil {
				return nil, err
			}

			accountActivities = append(accountActivities, accountActivitiesFromOperation...)
		}
	}
	return accountActivities, nil
}

func operationToAccountActivitySequence(sequence *model.Sequence, rawOperation *figmentclient.Operation, txHash string) ([]model.AccountActivitySeq, error) {
	var accountActivities []model.AccountActivitySeq

	marshaledData, err := json.Marshal(rawOperation.Details)
	if err != nil {
		return nil, err
	}

	switch rawOperation.Name {

	// Internal transfer
	case figmentclient.OperationTypeInternalTransfer:
		var event figmentclient.Transfer

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.From,
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeInternalTransferSent,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.To,
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeInternalTransferReceived,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

	//Election
	case figmentclient.OperationTypeValidatorGroupVoteCast:
		// vote() [ValidatorGroupVoteCast] => lockNonVoting->lockVotingPending
		var event contracts.ElectionValidatorGroupVoteCast

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupVoteCastSent,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Group.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupVoteCastReceived,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	case figmentclient.OperationTypeValidatorGroupVoteActivated:
		// activate() [ValidatorGroupVoteActivated] => lockVotingPending->lockVotingActive
		var event contracts.ElectionValidatorGroupVoteActivated

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupVoteActivatedSent,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Group.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupVoteActivatedReceived,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	case figmentclient.OperationTypeValidatorGroupPendingVoteRevoked:
		// revokePending() [ValidatorGroupPendingVoteRevoked] => lockVotingPending->lockNonVoting
		var event contracts.ElectionValidatorGroupPendingVoteRevoked

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupPendingVoteRevokedSent,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Group.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupPendingVoteRevokedReceived,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	case figmentclient.OperationTypeValidatorGroupActiveVoteRevoked:
		// revokeActive() [ValidatorGroupActiveVoteRevoked] => lockVotingActive->lockNonVoting
		var event contracts.ElectionValidatorGroupActiveVoteRevoked

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupActiveVoteRevokedSent,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Group.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            OperationTypeValidatorGroupActiveVoteRevokedReceived,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

	// Gold locked
	case figmentclient.OperationTypeGoldLocked:
		// lock() [GoldLocked + transfer] => main->lockNonVoting
		var event contracts.LockedGoldGoldLocked

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            rawOperation.Name,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	case figmentclient.OperationTypeGoldRelocked:
		// relock() [GoldRelocked] => lockPending->lockNonVoting
		var event contracts.LockedGoldGoldRelocked

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            rawOperation.Name,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	case figmentclient.OperationTypeGoldUnlocked:
		// unlock() [GoldUnlocked] => lockNonVoting->lockPending
		var event contracts.LockedGoldGoldUnlocked

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            rawOperation.Name,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	case figmentclient.OperationTypeGoldWithdrawn:
		// withdraw() [GoldWithdrawn + transfer] => lockPending->main
		var event contracts.LockedGoldGoldWithdrawn

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Account.String(),
			Amount:          types.NewQuantity(event.Value),
			Kind:            rawOperation.Name,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

	// Account
	case figmentclient.OperationTypeAccountSlashed:
		// slash() [AccountSlashed + transfer] => account:lockNonVoting -> beneficiary:lockNonVoting + governance:main
		var event contracts.LockedGoldAccountSlashed

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Slashed.String(),
			Amount:          types.NewQuantity(event.Penalty),
			Kind:            rawOperation.Name,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

	// Validators
	case figmentclient.OperationTypeValidatorEpochPaymentDistributed:
		var event contracts.ValidatorsValidatorEpochPaymentDistributed

		err := json.Unmarshal(marshaledData, &event)
		if err != nil {
			return nil, err
		}

		// For group
		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Group.String(),
			Amount:          types.NewQuantity(event.GroupPayment),
			Kind:            OperationTypeValidatorEpochPaymentDistributedForGroup,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})

		// For validator
		accountActivities = append(accountActivities, model.AccountActivitySeq{
			Sequence:        sequence,
			TransactionHash: txHash,
			Address:         event.Validator.String(),
			Amount:          types.NewQuantity(event.ValidatorPayment),
			Kind:            OperationTypeValidatorEpochPaymentDistributedForValidator,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	}

	return accountActivities, nil
}

func ToGovernanceActivitySequence(syncable *model.Syncable, parsedGovernanceLogs []*ParsedGovernanceLogs) ([]model.GovernanceActivitySeq, error) {
	var governanceActivities []model.GovernanceActivitySeq
	for _, log := range parsedGovernanceLogs {
		marshaledData, err := json.Marshal(log.Details)
		if err != nil {
			return nil, err
		}

		governanceActivities = append(governanceActivities, model.GovernanceActivitySeq{
			Sequence: &model.Sequence{
				Height: syncable.Height,
				Time:   *syncable.Time,
			},

			TransactionHash: log.TransactionHash,
			ProposalId:      log.ProposalId,
			Account:         log.Account,
			Kind:            log.Kind,
			Data:            types.Jsonb{RawMessage: marshaledData},
		})
	}
	return governanceActivities, nil
}
