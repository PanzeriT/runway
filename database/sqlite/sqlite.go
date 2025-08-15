package sqlite

import (
	"github.com/panzerit/runway/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteDriver struct{}

func (d *sqliteDriver) Connect(config []string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
	database.Register(&sqliteDriver{})
}
