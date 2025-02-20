// user_service.go

package services

import (
	"errors"

	"github.com/net22sky/telegram-bot/db/models"
	"github.com/net22sky/telegram-bot/db/repositories"

)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser создает нового пользователя или обновляет существующего.
func (s *UserService) CreateUser(telegramID int64, username, firstName string) (*models.User, error) {
	if telegramID <= 0 {
		return nil, errors.New("некорректный Telegram ID")
	}
	return s.userRepo.CreateUser(telegramID, username, firstName)
}

// GetUserByID получает пользователя по Telegram ID.
func (s *UserService) GetUserByID(telegramID int64) (*models.User, error) {
	if telegramID <= 0 {
		return nil, errors.New("некорректный Telegram ID")
	}
	return s.userRepo.GetUserByID(telegramID)
}


// UpdateUser обновляет данные пользователя.
func (s *UserService) UpdateUser(telegramID int64, updates map[string]interface{}) error {
    if telegramID <= 0 {
        return errors.New("некорректный Telegram ID")
    }
    return s.userRepo.UpdateUser(telegramID, updates)
}