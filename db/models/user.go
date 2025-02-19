package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TelegramID int64     `gorm:"unique;not null"`
	IsBot      bool      `gorm:"default:false"`
	Username   string    `gorm:"default:null"`
	FirstName  string    `gorm:"default:null"`
	LastName   string    `gorm:"default:null"`
	IsPremium  bool      `gorm:"default:false"`
	Language   string    `gorm:"default:'ru'"`
	Timezone   string    `gorm:"default:'UTC'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}