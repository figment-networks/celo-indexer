package store

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type BaseStore interface {
	Create(record interface{}) error
	Update(record interface{}) error
	Save(record interface{}) error
}
