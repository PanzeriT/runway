package service

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/panzerit/runway/model"
	"gorm.io/gorm"
)

type Service interface {
	GetModelNames() []string

	CreateRow(model string, dbObj any) error
	FindRows(model string, limit, offset int) (any, error)
	DeleteRow(model string, id uuid.UUID) error
}

type service struct {
	db     *gorm.DB
	models func() map[string]reflect.Type
}

func New(db *gorm.DB, models func() map[string]reflect.Type) Service {
	// TODO: this needs to be dynamic as well
	db.AutoMigrate(&model.User{}, &model.Role{})

	return &service{
		db:     db,
		models: models,
	}
}
