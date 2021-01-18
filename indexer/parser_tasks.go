package indexer

import (
	"context"
	"fmt"
	"github.com/celo-org/kliento/contracts"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
)

const (
	TaskNameGovernanceLogsParser = "GovernanceLogsParser"
)

var (
	_ pipeline.Task = (*governanceLogsParserTask)(nil)
)

//NewGovernanceLogsParserTask parses transaction logs to data about governance
func NewGovernanceLogsParserTask() *governanceLogsParserTask {
	return &governanceLogsParserTask{
		metricObserver: indexerTaskDuration.WithLabels(TaskNameGovernanceLogsParser),
	}
}

type governanceLogsParserTask struct {
	metricObserver metrics.Observer
}

type ParsedGovernanceLogs struct {
	ProposalId      uint64
	Account         string
	TransactionHash string
	Kind            string
	Details         interface{}
}

func (t *governanceLogsParserTask) GetName() string {
	return TaskNameGovernanceLogsParser
}

func (t *governanceLogsParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageParser, t.GetName(), payload.CurrentHeight))

	rawTransactions := payload.RawTransactions

	var parsedGovernanceLogsData []*ParsedGovernanceLogs

	for _, rawTransaction := range rawTransactions {
		for _, rawOperation := range rawTransaction.Operations {

			parsedLog := getGovernanceLogData(rawOperation)
			if parsedLog != nil {
				parsedLog.TransactionHash = rawTransaction.Hash

				parsedGovernanceLogsData = append(parsedGovernanceLogsData, parsedLog)
			}
		}
	}

	payload.ParsedGovernanceLogs = parsedGovernanceLogsData

	return nil
}

func getGovernanceLogData(rawOperation *figmentclient.Operation) *ParsedGovernanceLogs {
	switch rawOperation.Name {

	case figmentclient.OperationTypeProposalVoted:
		event := rawOperation.Details.(*contracts.GovernanceProposalVoted)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Account:    event.Account.String(),
			Kind:       rawOperation.Name,
			Details:    event,
		}

	case figmentclient.OperationTypeProposalUpvoted:
		event := rawOperation.Details.(*contracts.GovernanceProposalUpvoted)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Account:    event.Account.String(),
			Kind:       rawOperation.Name,
			Details:    event,
		}

	case figmentclient.OperationTypeProposalApproved:
		event := rawOperation.Details.(*contracts.GovernanceProposalApproved)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Kind:       rawOperation.Name,
			Details:    event,
		}

	case figmentclient.OperationTypeProposalExecuted:
		event := rawOperation.Details.(*contracts.GovernanceProposalExecuted)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Kind:       rawOperation.Name,
			Details:    event,
		}

	case figmentclient.OperationTypeProposalDequeued:
		event := rawOperation.Details.(*contracts.GovernanceProposalDequeued)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Kind:       rawOperation.Name,
			Details:    event,
		}

	case figmentclient.OperationTypeProposalQueued:
		event := rawOperation.Details.(*contracts.GovernanceProposalQueued)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Account:    event.Proposer.String(),
			Kind:       rawOperation.Name,
			Details:    event,
		}

	case figmentclient.OperationTypeProposalExpired:
		event := rawOperation.Details.(*contracts.GovernanceProposalExpired)
		return &ParsedGovernanceLogs{
			ProposalId: event.ProposalId.Uint64(),
			Kind:       rawOperation.Name,
			Details:    event,
		}
	}

	return nil
}
