package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Store{Db: conn}, nil
}

// Store handles all database operations
type Store struct {
	Db *gorm.DB
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.Db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.Db.Close()
}
