package store

import (
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewAccountActivitySeqStore(db *gorm.DB) *AccountActivitySeqStore {
	return &AccountActivitySeqStore{scoped(db, model.AccountActivitySeq{})}
}

// AccountActivitySeqStore handles operations on account activities
type AccountActivitySeqStore struct {
	baseStore
}

// CreateIfNotExists creates the account activity if it does not exist
func (s AccountActivitySeqStore) CreateIfNotExists(accountActivity *model.AccountActivitySeq) error {
	_, err := s.FindByHeight(accountActivity.Height)
	if isNotFound(err) {
		return s.Create(accountActivity)
	}
	return nil
}

// FindByHeightAndAddress finds account activities by height and address
func (s AccountActivitySeqStore) FindByHeightAndAddress(height int64, address string) ([]model.AccountActivitySeq, error) {
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
func (s AccountActivitySeqStore) FindByHeight(h int64) ([]model.AccountActivitySeq, error) {
	var result []model.AccountActivitySeq

	err := s.db.
		Where("height = ?", h).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindMostRecent finds most recent account activity sequence
func (s *AccountActivitySeqStore) FindMostRecent() (*model.AccountActivitySeq, error) {
	accountActivitySeq := &model.AccountActivitySeq{}
	if err := findMostRecent(s.db, "time", accountActivitySeq); err != nil {
		return nil, err
	}
	return accountActivitySeq, nil
}

// FindLastByAddress finds last account activity sequences for given address
func (s AccountActivitySeqStore) FindLastByAddress(address string, limit int64) ([]model.AccountActivitySeq, error) {
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
func (s AccountActivitySeqStore) FindLastByAddressAndKind(address string, kind string, limit int64) ([]model.AccountActivitySeq, error) {
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
func (s *AccountActivitySeqStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
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
func (s *AccountActivitySeqStore) DeleteForHeight(h int64) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("height = ?", h).
		Delete(&model.AccountActivitySeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}