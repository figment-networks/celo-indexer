package store

import (
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewGovernanceActivitySeqStore(db *gorm.DB) *GovernanceActivitySeqStore {
	return &GovernanceActivitySeqStore{scoped(db, model.GovernanceActivitySeq{})}
}

// GovernanceActivitySeqStore handles operations on governance activities
type GovernanceActivitySeqStore struct {
	baseStore
}

// CreateIfNotExists creates the governance activity if it does not exist
func (s GovernanceActivitySeqStore) CreateIfNotExists(governanceActivity *model.GovernanceActivitySeq) error {
	_, err := s.FindByHeight(governanceActivity.Height)
	if isNotFound(err) {
		return s.Create(governanceActivity)
	}
	return nil
}

// FindByHeightAndProposalId finds governance activities by height and proposal Id
func (s GovernanceActivitySeqStore) FindByHeightAndProposalId(height int64, proposalId uint64) ([]model.GovernanceActivitySeq, error) {
	q := model.GovernanceActivitySeq{
		ProposalId: proposalId,
	}
	var result []model.GovernanceActivitySeq

	err := s.db.
		Where(&q).
		Where("height = ?", height).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindByHeight finds governance activity sequences by height
func (s GovernanceActivitySeqStore) FindByHeight(h int64) ([]model.GovernanceActivitySeq, error) {
	var result []model.GovernanceActivitySeq

	err := s.db.
		Where("height = ?", h).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindByProposalId finds governance activities by proposal Id
func (s GovernanceActivitySeqStore) FindByProposalId(proposalId uint64, limit int64, cursor *int64) ([]model.GovernanceActivitySeq, *int64, error) {
	q := model.GovernanceActivitySeq{
		ProposalId: proposalId,
	}
	var result []model.GovernanceActivitySeq

	tx := s.db.
		Where(&q).
		Order("id DESC")

	if cursor != nil {
		tx = tx.Where("id < ?", cursor)
	}

	tx = tx.
		Limit(limit).
		Find(&result)

	if tx.Error != nil {
		return nil, nil, checkErr(tx.Error)
	}

	var nextCursor int64
	if len(result) > 0 {
		nextCursor = int64(result[len(result) - 1].ID)
	}

	return result, &nextCursor, nil
}

// FindMostRecent finds most recent governance activity sequence
func (s *GovernanceActivitySeqStore) FindMostRecent() (*model.GovernanceActivitySeq, error) {
	governanceActivitySeq := &model.GovernanceActivitySeq{}
	if err := findMostRecent(s.db, "time", governanceActivitySeq); err != nil {
		return nil, err
	}
	return governanceActivitySeq, nil
}

// FindLastByProposalId finds last governance activity sequences for given proposal Id
func (s GovernanceActivitySeqStore) FindLastByProposalId(proposalId uint64, limit int64) ([]model.GovernanceActivitySeq, error) {
	q := model.GovernanceActivitySeq{
		ProposalId: proposalId,
	}
	var result []model.GovernanceActivitySeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindLastByProposalIdAndKind finds last governance activity sequences for given proposal Id and kind
func (s GovernanceActivitySeqStore) FindLastByProposalIdAndKind(proposalId uint64, kind string, limit int64) ([]model.GovernanceActivitySeq, error) {
	q := model.GovernanceActivitySeq{
		ProposalId: proposalId,
		Kind:       kind,
	}
	var result []model.GovernanceActivitySeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// DeleteOlderThan deletes governance activity sequence older than given threshold
func (s *GovernanceActivitySeqStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.GovernanceActivitySeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

// DeleteForHeight deletes governance activity sequence for given height
func (s *GovernanceActivitySeqStore) DeleteForHeight(h int64) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("height = ?", h).
		Delete(&model.GovernanceActivitySeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}
