package models

import "time"

type Note struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Text      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"constraint:OnUpdate:CASCADE;OnDelete:CASCADE;"`
}