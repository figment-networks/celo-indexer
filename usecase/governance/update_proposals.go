package governance

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/client/theceloclient"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/pkg/errors"
	"strconv"
)

var (
	ErrSomeProposalsNotUpdated = errors.New("error occurred while updating some proposals.")
)

type updateProposalsUseCase struct {
	proposalAggDb store.ProposalAgg
	client        theceloclient.Client
}

func NewUpdateProposalsUseCase(c theceloclient.Client, proposalAggDb store.ProposalAgg) *updateProposalsUseCase {
	return &updateProposalsUseCase{
		client:        c,
		proposalAggDb: proposalAggDb,
	}
}

func (uc *updateProposalsUseCase) Execute(ctx context.Context) error {
	logger.Info(fmt.Sprintf("running update proposals use case [handler=cmd]"))

	persistedProposals, _, err := uc.proposalAggDb.All(0, nil)
	if err != nil {
		return err
	}

	sourceProposals, err := uc.client.GetAllProposals()
	if err != nil {
		return err
	}

	errorsCount := 0
	for _, persistedProposal := range persistedProposals {
		proposalId := strconv.FormatUint(persistedProposal.ProposalId, 10)
		sourceProposal, ok := sourceProposals.Items[proposalId]
		if ok && persistedProposal.DescriptionUrl == "" {
			persistedProposal.DescriptionUrl = sourceProposal.DescriptionUrl

			if err = uc.proposalAggDb.Save(&persistedProposal); err != nil {
				errorsCount++
			}
		}
	}

	if errorsCount > 0 {
		return ErrSomeProposalsNotUpdated
	}

	return nil
}
