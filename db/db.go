package db

import (
	"fmt"
	"github.com/net22sky/telegram-bot/db/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitDB инициализирует подключение к базе данных и возвращает экземпляр *gorm.DB.
func InitDB(dsn string, migrate bool) (*gorm.DB, error) {
    // Подключаемся к базе данных
    dbInstance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
    }

    log.Println("Подключено к базе данных через GORM")

    // Выполняем миграции, если включено
    if migrate {
        err = migrateDB(dbInstance)
        if err != nil {
            return nil, fmt.Errorf("ошибка при миграции таблиц: %w", err)
        }
    }

    return dbInstance, nil
}

// migrateDB выполняет миграции для всех моделей.
func migrateDB(db *gorm.DB) error {
    models := []interface{}{
        &models.User{},
        &models.Note{},
        &models.PollAnswer{},
    }

    for _, model := range models {
        err := db.AutoMigrate(model)
        if err != nil {
            return fmt.Errorf("ошибка при миграции модели %T: %w", model, err)
        }
        log.Printf("Миграция выполнена для модели %T", model)
    }

    return nil
}