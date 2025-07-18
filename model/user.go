package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey;unique"`
	Email         string    `gorm:"unique"`
	FirstName     string
	LastName      string
	IsAdmin       bool
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	DeactivatedAt time.Time
}
