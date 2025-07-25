package service

import (
	"github.com/google/uuid"
	"github.com/panzerit/runway/model"

	"gorm.io/gorm"
)

type Service interface {
	GetModelNames() []string

	CreateRow(dbObj any) error
	CreateRowFromJSON(model string, jsonData []byte) error
	FindRows(model string, limit, offset int) (any, error)
	DeleteRow(model string, id uuid.UUID) error
}

type service struct {
	db     *gorm.DB
	models func() map[string]model.Model
}

func New(db *gorm.DB, models func() map[string]model.Model) Service {
	// TODO: this needs to be dynamic as well
	db.AutoMigrate(&model.User{}, &model.Role{})

	return &service{
		db:     db,
		models: models,
	}
}
