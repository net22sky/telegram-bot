package models

import "time"

type PollAnswer struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	PollID    string    `gorm:"not null"`
	OptionIDs []int     `gorm:"serializer:json"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"constraint:OnUpdate:CASCADE;OnDelete:CASCADE;"`
}