package models


import "time"

// BirthdayReminder представляет напоминание о дне рождения.
type BirthdayReminder struct {
    UserID        int64        `gorm:"index"`
    Name          string
    Day           int          `gorm:"index"`
    Month         int          `gorm:"index"`
    CreatedAt     time.Time    `gorm:"autoCreateTime"`
    User          User         `gorm:"constraint:OnUpdate:CASCADE;OnDelete:CASCADE;"`
    ReminderTypeID int64       `gorm:"index"` // Индекс для внешнего ключа
    ReminderType   ReminderType `gorm:"foreignKey:ReminderTypeID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}

// ReminderType представляет тип напоминания.
type ReminderType struct {
    ID   int64  `gorm:"primaryKey"` // Уникальный идентификатор типа
    Name string // Название типа (например, "День рождения")
}