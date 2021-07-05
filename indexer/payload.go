package indexer

import (
	"fmt"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/indexing-engine/datalake"
	"github.com/figment-networks/indexing-engine/pipeline"
)

var (
	_ pipeline.PayloadFactory = (*payloadFactory)(nil)
	_ pipeline.Payload        = (*payload)(nil)
)

func NewPayloadFactory(dl *datalake.DataLake) *payloadFactory {
	return &payloadFactory{dl: dl}
}

type payloadFactory struct {
	dl *datalake.DataLake
}

func (pf *payloadFactory) GetPayload(currentHeight int64) pipeline.Payload {
	return &payload{
		CurrentHeight: currentHeight,
		DataLake:      pf.dl,
	}
}

type payload struct {
	CurrentHeight int64

	DataLake *datalake.DataLake

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

func (p *payload) Store(name string, obj interface{}) error {
	res, err := datalake.NewJSONResource(obj)
	if err != nil {
		return fmt.Errorf("cannot store %s in data lake [height=%d]: %v",
			name, p.CurrentHeight, err)
	}

	return p.DataLake.StoreResourceAtHeight(res, name, p.CurrentHeight)
}

func (p *payload) Retrieve(name string, obj interface{}) error {
	res, err := p.DataLake.RetrieveResourceAtHeight(name, p.CurrentHeight)
	if err != nil {
		return fmt.Errorf("cannot retrieve %s from data lake [height=%d]: %v",
			name, p.CurrentHeight, err)
	}

	return res.ScanJSON(obj)
}
