package store

import (
	"github.com/figment-networks/celo-indexer/model"
	"time"
)

type AccountActivitySeq interface {
	BulkUpsert(records []model.AccountActivitySeq) error
	FindByHeightAndAddress(height int64, address string) ([]model.AccountActivitySeq, error)
	FindByHeight(h int64) ([]model.AccountActivitySeq, error)
	FindMostRecent() (*model.AccountActivitySeq, error)
	FindLastByAddress(address string, limit int64) ([]model.AccountActivitySeq, error)
	FindLastByAddressAndKind(address string, kind string, limit int64) ([]model.AccountActivitySeq, error)
	DeleteOlderThan(purgeThreshold time.Time) (*int64, error)
	DeleteForHeight(h int64) (*int64, error)
}

