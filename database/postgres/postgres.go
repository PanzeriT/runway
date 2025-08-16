package postgres

import (
	"fmt"

	"github.com/panzerit/runway/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds the configuration for connecting to a PostgreSQL database.
//
// If you don't want to add a TimeZone, you can leave it empty. Some PostgreSQL Servers
// do not support the TimeZone parameter, e.g. Google Cloud SQL.
type Config struct {
	Host     string
	User     string
	Password string
	Port     string
	Name     string
	SSLMode  string
	TimeZone string
}

func (c Config) Type() string {
	return "postgres"
}

// NewConfig creates a postgres config with sensible defaults
func NewConfig() Config {
	return Config{
		Host:     "localhost",
		User:     "runway",
		Password: "password",
		Port:     "5432",
		Name:     "runway",
		SSLMode:  "disable",
		TimeZone: "Europe/Zurich",
	}
}

type driver struct{}

func (d *driver) Connect(config database.Config) (*gorm.DB, error) {
	pgConfig, ok := config.(Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type for postgres driver")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		pgConfig.Host,
		pgConfig.User,
		pgConfig.Password,
		pgConfig.Name,
		pgConfig.Port,
		pgConfig.SSLMode,
	)

	if pgConfig.TimeZone != "" {
		dsn += fmt.Sprintf(" TimeZone=%s", pgConfig.TimeZone)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}

func init() {
	database.Register("postgres", &driver{})
}
