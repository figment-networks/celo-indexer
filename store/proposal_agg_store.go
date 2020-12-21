package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewProposalAggStore(db *gorm.DB) *ProposalAggStore {
	return &ProposalAggStore{scoped(db, model.ProposalAgg{})}
}

// ProposalAggStore handles operations on proposals
type ProposalAggStore struct {
	baseStore
}

// CreateOrUpdate creates a new proposal or updates an existing one
func (s ProposalAggStore) CreateOrUpdate(val *model.ProposalAgg) error {
	existing, err := s.FindByProposalId(val.ProposalId)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}
	return s.Update(existing)
}

// FindBy returns an proposal for a matching attribute
func (s ProposalAggStore) FindBy(key string, value interface{}) (*model.ProposalAgg, error) {
	result := &model.ProposalAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns an proposal for the ID
func (s ProposalAggStore) FindByID(id int64) (*model.ProposalAgg, error) {
	return s.FindBy("id", id)
}

// FindByProposalId return proposal by proposal Id
func (s *ProposalAggStore) FindByProposalId(proposalId uint64) (*model.ProposalAgg, error) {
	return s.FindBy("proposal_id", proposalId)
}

// All returns all proposals
func (s ProposalAggStore) All(limit int64, cursor *int64) ([]model.ProposalAgg, *int64, error) {
	var result []model.ProposalAgg

	tx := s.db.
		Order("id DESC")

	if cursor != nil {
		tx = tx.Where("id < ?", cursor)
	}

	if limit > 0 {
		tx = tx.Limit(limit)
	}

	tx = tx.Find(&result)

	if tx.Error != nil {
		return nil, nil, checkErr(tx.Error)
	}

	var nextCursor int64
	if len(result) > 0 {
		nextCursor = int64(result[len(result) - 1].ID)
	}

	return result, &nextCursor, nil
}
