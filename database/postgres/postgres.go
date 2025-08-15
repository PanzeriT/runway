package postgres

import (
	"fmt"
	"reflect"

	"github.com/panzerit/runway/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type driver struct{}

type config struct {
	host     string
	user     string
	password string
	port     string
	name     string
	sslMode  string
}

const (
	host     database.ConfigKey = "host"
	user     database.ConfigKey = "user"
	password database.ConfigKey = "password"
	port     database.ConfigKey = "port"
	name     database.ConfigKey = "name"
	sslMode  database.ConfigKey = "sslmode"
)

var defaults = map[database.ConfigKey]string{
	host:     "localhost",
	user:     "runway",
	password: "password",
	port:     "5432",
	name:     "runway",
	sslMode:  "disable",
}

func (d *driver) validateKeys(keyValueConfig any) config {
	c := config{
		host:     defaults[host],
		user:     defaults[user],
		password: defaults[password],
		port:     defaults[port],
		name:     defaults[name],
		sslMode:  defaults[sslMode],
	}

	v := reflect.ValueOf(c).Elem()

	configArray := keyValueConfig.([]string)
	for i := 0; i < len(configArray); i += 2 {
		name := configArray[i]
		value := configArray[i+1]
		f := v.FieldByName(name)
		if !f.IsValid() {
			panic("No such field: " + name)
		}
		if !f.CanSet() {
			panic("Cannot set field: " + name)
		}
		f.Set(reflect.ValueOf(value))
	}

	return c
}

func (d *driver) Connect(config []string) (*gorm.DB, error) {
	c := d.validateKeys(config)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.host,
		c.user,
		c.password,
		c.name,
		c.port,
		c.sslMode,
		"",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
	database.Register(&driver{})
}
