package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"time"
)

type ProposalAgg interface {
	Create(*model.ProposalAgg) error
	Save(*model.ProposalAgg) error
	CreateOrUpdate(val *model.ProposalAgg) error
	FindBy(key string, value interface{}) (*model.ProposalAgg, error)
	FindByID(id int64) (*model.ProposalAgg, error)
	FindByProposalId(proposalId uint64) (*model.ProposalAgg, error)
	All(limit int64, cursor *int64) ([]model.ProposalAgg, *int64, error)
}

type GovernanceActivitySeq interface {
	BulkUpsert(records []model.GovernanceActivitySeq) error
	CreateIfNotExists(governanceActivity *model.GovernanceActivitySeq) error
	FindByHeightAndProposalId(height int64, proposalId uint64) ([]model.GovernanceActivitySeq, error)
	FindByHeight(h int64) ([]model.GovernanceActivitySeq, error)
	FindByProposalId(proposalId uint64, limit int64, cursor *int64) ([]model.GovernanceActivitySeq, *int64, error)
	FindMostRecent() (*model.GovernanceActivitySeq, error)
	FindLastByProposalId(proposalId uint64, limit int64) ([]model.GovernanceActivitySeq, error)
	FindLastByProposalIdAndKind(proposalId uint64, kind string, limit int64) ([]model.GovernanceActivitySeq, error)
	DeleteOlderThan(purgeThreshold time.Time) (*int64, error)
	DeleteForHeight(h int64) (*int64, error)
}
