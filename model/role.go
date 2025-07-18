package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID            uuid.UUID `gorm:"primaryKey;unique" runway:"hidden"`
	Name          string    `gorm:"unique"`
	CreatedAt     time.Time `gorm:"autrCreateTime" runway:"hidden"`
	DeactivatedAt time.Time `runway:"hidden"`
}
