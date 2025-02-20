package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TelegramID int64     `gorm:"unique;not null"` // ID пользователя в Telegram
	IsBot      bool      `gorm:"default:false"`   // Является ли пользователь ботом
	Username   string    `gorm:"default:null"`    // Имя пользователя
	FirstName  string    `gorm:"default:null"`    // Имя
	LastName   string    `gorm:"default:null"`    // Фамилия
	IsPremium  bool      `gorm:"default:false"`   // Премиум-статус
	Language   string    `gorm:"default:'ru'"`    // Язык пользователя
	Timezone   string    `gorm:"default:'UTC'"`   // Часовой пояс
	CreatedAt  time.Time `gorm:"autoCreateTime"`  // Время создания записи
}
