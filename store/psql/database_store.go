package psql

import (
	"github.com/figment-networks/celo-indexer/store"
	"github.com/jinzhu/gorm"
)

var _ store.Database = (*Database)(nil)

func NewDatabaseStore(db *gorm.DB) *Database {
	return &Database{
		db: db,
	}
}

// Database handles operations on blocks
type Database struct {
	db *gorm.DB
}

// FindSummary Gets average block times for interval
func (s *Database) GetTotalSize() (*store.GetTotalSizeResult, error) {
	query := "SELECT pg_database_size(current_database()) as size"

	var result store.GetTotalSizeResult
	err := s.db.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
