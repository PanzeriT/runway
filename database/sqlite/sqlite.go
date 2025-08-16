package sqlite

import (
	"fmt"

	"github.com/panzerit/runway/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Path string
}

func (c Config) Type() string {
	return "sqlite"
}

// NewConfig creates a sqlite config with sensible defaults
func NewConfig() Config {
	return Config{
		Path: "app.db",
	}
}

type driver struct{}

func (d *driver) Connect(config database.Config) (*gorm.DB, error) {
	sqliteConfig, ok := config.(Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type for sqlite driver")
	}

	db, err := gorm.Open(sqlite.Open(sqliteConfig.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
	}

	return db, nil
}

func init() {
	database.Register("sqlite", &driver{})
}
