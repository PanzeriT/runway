// database.go
package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Driver interface {
	Connect(config Config) (*gorm.DB, error)
}

type Config interface {
	Type() string
}

var drivers = make(map[string]Driver)

func Register(driverType string, driver Driver) {
	if _, exists := drivers[driverType]; exists {
		panic(fmt.Sprintf("driver %s already registered", driverType))
	}
	drivers[driverType] = driver
}

func Init(config Config) (*gorm.DB, error) {
	driver, exists := drivers[config.Type()]
	if !exists {
		return nil, fmt.Errorf("no driver registered for type: %s", config.Type())
	}
	return driver.Connect(config)
}
