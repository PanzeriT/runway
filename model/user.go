package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey;unique" runway:"hidden" json:"id"`
	Email         string    `gorm:"unique" runway:"create" json:"email"`
	FirstName     string    `runway:"create" json:"first_name"`
	LastName      string    `runway:"create" json:"last_name"`
	IsAdmin       bool      `json:"is_admin"`
	CreatedAt     time.Time `gorm:"autoCreateTime" runway:"hidden" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" runway:"hidden" json:"updated_at"`
	DeactivatedAt time.Time `runway:"hidden" json:"deactivated_at"`
}
