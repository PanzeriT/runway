package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey;unique" runway:"hidden"`
	Email         string    `gorm:"unique"`
	FirstName     string
	LastName      string
	IsAdmin       bool
	CreatedAt     time.Time `gorm:"autoCreateTime" runway:"hidden"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" runway:"hidden"`
	DeactivatedAt time.Time `runway:"hidden"`
}
