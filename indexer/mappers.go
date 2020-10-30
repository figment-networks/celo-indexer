package indexer

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/pkg/errors"
)

var (
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

func ToValidatorSequence(syncable *model.Syncable, rawValidatorPerformance []*figmentclient.Validator) ([]model.ValidatorSeq, error) {
	var validators []model.ValidatorSeq
	for _, rawValidator := range rawValidatorPerformance {
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

func ToValidatorGroupSequence(syncable *model.Syncable, rawValidatorGroups []*figmentclient.ValidatorGroup) ([]model.ValidatorGroupSeq, error) {
	var validators []model.ValidatorGroupSeq
	for _, rawValidator := range rawValidatorGroups {
		e := model.ValidatorGroupSeq{
			Sequence: &model.Sequence{
				Height: syncable.Height,
				Time:   *syncable.Time,
			},

			Address:         rawValidator.Address,
			Commission:      types.NewQuantity(rawValidator.Commission),
			ActiveVotes:     types.NewQuantity(rawValidator.ActiveVotes),
			ActiveVoteUnits: types.NewQuantity(rawValidator.ActiveVotesUnits),
			PendingVotes:    types.NewQuantity(rawValidator.PendingVotes),
			MembersCount:    len(rawValidator.Members),
		}

		if !e.Valid() {
			return nil, ErrValidatorGroupSequenceNotValid
		}

		validators = append(validators, e)
	}
	return validators, nil
}
