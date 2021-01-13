package indexer

import (
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/indexing-engine/pipeline"
)

var (
	_ pipeline.PayloadFactory = (*payloadFactory)(nil)
	_ pipeline.Payload        = (*payload)(nil)
)

func NewPayloadFactory() *payloadFactory {
	return &payloadFactory{}
}

type payloadFactory struct{}

func (pf *payloadFactory) GetPayload(currentHeight int64) pipeline.Payload {
	return &payload{
		CurrentHeight: currentHeight,
	}
}

type payload struct {
	CurrentHeight int64

	// Fetcher stage
	HeightMeta         HeightMeta
	RawBlock           *figmentclient.Block
	RawValidators      []*figmentclient.Validator
	RawValidatorGroups []*figmentclient.ValidatorGroup
	RawTransactions    []*figmentclient.Transaction

	// Syncer stage
	Syncable *model.Syncable

	// Parser stage
	ParsedGovernanceLogs []*ParsedGovernanceLogs

	// Aggregator stage
	NewValidatorAggregates          []model.ValidatorAgg
	UpdatedValidatorAggregates      []model.ValidatorAgg
	NewValidatorGroupAggregates     []model.ValidatorGroupAgg
	UpdatedValidatorGroupAggregates []model.ValidatorGroupAgg
	NewProposalAggregates           []model.ProposalAgg
	UpdatedProposalAggregates       []model.ProposalAgg

	// Sequencer stage
	NewBlockSequence            *model.BlockSeq
	UpdatedBlockSequence        *model.BlockSeq
	ValidatorSequences          []model.ValidatorSeq
	ValidatorGroupSequences     []model.ValidatorGroupSeq
	AccountActivitySequences    []model.AccountActivitySeq
	GovernanceActivitySequences []model.GovernanceActivitySeq

	// Analyzer
	SystemEvents []model.SystemEvent
}

func (p *payload) MarkAsProcessed() {}
