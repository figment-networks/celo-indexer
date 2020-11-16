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

	// Aggregator stage
	NewValidatorAggregates          []model.ValidatorAgg
	UpdatedValidatorAggregates      []model.ValidatorAgg
	NewValidatorGroupAggregates     []model.ValidatorGroupAgg
	UpdatedValidatorGroupAggregates []model.ValidatorGroupAgg

	// Sequencer stage
	NewBlockSequence               *model.BlockSeq
	UpdatedBlockSequence           *model.BlockSeq
	NewValidatorSequences          []model.ValidatorSeq
	UpdatedValidatorSequences      []model.ValidatorSeq
	NewValidatorGroupSequences     []model.ValidatorGroupSeq
	UpdatedValidatorGroupSequences []model.ValidatorGroupSeq

	AccountActivitySequences []model.AccountActivitySeq

	// Analyzer
	SystemEvents []*model.SystemEvent
}

func (p *payload) MarkAsProcessed() {}
