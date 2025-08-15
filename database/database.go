package database

import (
	"gorm.io/gorm"
)

type Driver interface {
	Connect(config []string) (*gorm.DB, error)
}

type (
	Config    []string
	ConfigKey string
)

var driver Driver

func Register(d Driver) {
	if driver != nil {
		panic("driver already registered")
	}
	driver = d
}

func Init(config ...string) (*gorm.DB, error) {
	if driver == nil {
		panic("no driver registered")
	}

	if len(config)%2 != 0 {
		panic("config needs to be must be key-value pairs")
	}

	return driver.Connect(config)
}
