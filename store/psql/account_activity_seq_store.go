package psql

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/indexing-engine/store/bulk"
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
)

var _ store.AccountActivitySeq = (*AccountActivitySeq)(nil)

func NewAccountActivitySeqStore(db *gorm.DB) *AccountActivitySeq {
	return &AccountActivitySeq{scoped(db, model.AccountActivitySeq{})}
}

// AccountActivitySeq handles operations on account activities
type AccountActivitySeq struct {
	baseStore
}

// BulkUpsert insert account activity sequences in bulk
func (s AccountActivitySeq) BulkUpsert(records []model.AccountActivitySeq) error {
	var err error

	for i := 0; i < len(records); i += batchSize {
		j := i + batchSize
		if j > len(records) {
			j = len(records)
		}

		err = s.baseStore.BulkUpsert(bulkInsertAccountActivitySeqs, j-i, func(k int) bulk.Row {
			r := records[i+k]
			return bulk.Row{
				r.Height,
				r.Time,
				r.TransactionHash,
				r.Address,
				r.Amount.String(),
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

// FindByHeightAndAddress finds account activities by height and address
func (s AccountActivitySeq) FindByHeightAndAddress(height int64, address string) ([]model.AccountActivitySeq, error) {
	q := model.AccountActivitySeq{
		Address: address,
	}
	var result []model.AccountActivitySeq

	err := s.db.
		Where(&q).
		Where("height = ?", height).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindByHeight finds account activity sequences by height
func (s AccountActivitySeq) FindByHeight(h int64) ([]model.AccountActivitySeq, error) {
	var result []model.AccountActivitySeq

	err := s.db.
		Where("height = ?", h).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindMostRecent finds most recent account activity sequence
func (s *AccountActivitySeq) FindMostRecent() (*model.AccountActivitySeq, error) {
	accountActivitySeq := &model.AccountActivitySeq{}
	if err := findMostRecent(s.db, "time", accountActivitySeq); err != nil {
		return nil, err
	}
	return accountActivitySeq, nil
}

// FindLastByAddress finds last account activity sequences for given address
func (s AccountActivitySeq) FindLastByAddress(address string, limit int64) ([]model.AccountActivitySeq, error) {
	q := model.AccountActivitySeq{
		Address: address,
	}
	var result []model.AccountActivitySeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindLastByAddressAndKind finds last account activity sequences for given address and kind
func (s AccountActivitySeq) FindLastByAddressAndKind(address string, kind string, limit int64) ([]model.AccountActivitySeq, error) {
	q := model.AccountActivitySeq{
		Address: address,
		Kind:    kind,
	}
	var result []model.AccountActivitySeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// DeleteOlderThan deletes account activity sequence older than given threshold
func (s *AccountActivitySeq) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.AccountActivitySeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

// DeleteForHeight deletes account activity sequence for given height
func (s *AccountActivitySeq) DeleteForHeight(h int64) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("height = ?", h).
		Delete(&model.AccountActivitySeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}