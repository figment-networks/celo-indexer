package psql

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
)

var _ store.GovernanceActivitySeq = (*GovernanceActivitySeq)(nil)

func NewGovernanceActivitySeqStore(db *gorm.DB) *GovernanceActivitySeq {
	return &GovernanceActivitySeq{scoped(db, model.GovernanceActivitySeq{})}
}

// GovernanceActivitySeq handles operations on governance activities
type GovernanceActivitySeq struct {
	baseStore
}

// BulkUpsert insert validator sequences in bulk
func (s GovernanceActivitySeq) BulkUpsert(records []model.GovernanceActivitySeq) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertGovernanceActivitySeqs, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.Height,
				r.Time,
				r.ProposalId,
				r.Account,
				r.TransactionHash,
				r.Kind,
				r.Data,
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateIfNotExists creates the governance activity if it does not exist
func (s GovernanceActivitySeq) CreateIfNotExists(governanceActivity *model.GovernanceActivitySeq) error {
	_, err := s.FindByHeight(governanceActivity.Height)
	if isNotFound(err) {
		return s.Create(governanceActivity)
	}
	return nil
}

// FindByHeightAndProposalId finds governance activities by height and proposal Id
func (s GovernanceActivitySeq) FindByHeightAndProposalId(height int64, proposalId uint64) ([]model.GovernanceActivitySeq, error) {
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
func (s GovernanceActivitySeq) FindByHeight(h int64) ([]model.GovernanceActivitySeq, error) {
	var result []model.GovernanceActivitySeq

	err := s.db.
		Where("height = ?", h).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindByProposalId finds governance activities by proposal Id
func (s GovernanceActivitySeq) FindByProposalId(proposalId uint64, limit int64, cursor *int64) ([]model.GovernanceActivitySeq, *int64, error) {
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
func (s *GovernanceActivitySeq) FindMostRecent() (*model.GovernanceActivitySeq, error) {
	governanceActivitySeq := &model.GovernanceActivitySeq{}
	if err := findMostRecent(s.db, "time", governanceActivitySeq); err != nil {
		return nil, err
	}
	return governanceActivitySeq, nil
}

// FindLastByProposalId finds last governance activity sequences for given proposal Id
func (s GovernanceActivitySeq) FindLastByProposalId(proposalId uint64, limit int64) ([]model.GovernanceActivitySeq, error) {
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
func (s GovernanceActivitySeq) FindLastByProposalIdAndKind(proposalId uint64, kind string, limit int64) ([]model.GovernanceActivitySeq, error) {
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
func (s *GovernanceActivitySeq) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
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
func (s *GovernanceActivitySeq) DeleteForHeight(h int64) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("height = ?", h).
		Delete(&model.GovernanceActivitySeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}
