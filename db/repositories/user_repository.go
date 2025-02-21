// user_repository.go

package repositories

import (
	"fmt"
	"log"

	"github.com/net22sky/telegram-bot/db/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя или обновляет существующего.
func (r *UserRepository) CreateUser(telegramID int64, username, firstName string) (*models.User, error) {
	user := models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
	}
	result := r.db.FirstOrCreate(&user, models.User{TelegramID: telegramID})
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", result.Error)
	}

	// Обновляем поля, если пользователь уже существует
	if result.RowsAffected == 0 {
		r.db.Model(&user).Updates(models.User{
			Username:  username,
			FirstName: firstName,
		})
		log.Printf("Пользователь %d успешно обновлен: %+v", telegramID, user)
	} else {
		log.Printf("Пользователь %d успешно создан: %+v", telegramID, user)
	}

	return &user, nil
}

// GetUserByID получает пользователя по Telegram ID.
func (r *UserRepository) GetUserByID(telegramID int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "telegram_id = ?", telegramID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Пользователь %d не найден", telegramID)
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", result.Error)
	}
	log.Printf("Пользователь %d успешно загружен: %+v", telegramID, user)
	return &user, nil
}

// UpdateUser обновляет данные пользователя.
func (r *UserRepository) UpdateUser(telegramID int64, updates map[string]interface{}) error {
	var user models.User
	result := r.db.Model(&user).Where("telegram_id = ?", telegramID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("пользователь с Telegram ID %d не найден", telegramID)
	}
	log.Printf("Пользователь %d успешно обновлен: %+v", telegramID, updates)
	return nil
}

// UpdateUserLanguage обновляет язык пользователя.
func (r *UserRepository) UpdateUserLanguage(telegramID int64, language string) error {
	return r.db.Model(&models.User{}).Where("telegram_id = ?", telegramID).Update("language", language).Error
}
