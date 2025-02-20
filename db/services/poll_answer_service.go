package services

import (
	"github.com/net22sky/telegram-bot/db/repositories"
	"gorm.io/gorm"
)

// PollAnswerService предоставляет методы для работы с ответами на опросы
type PollAnswerService struct {
	repo *repositories.PollAnswerRepository
}

// NewPollAnswerService создает новый экземпляр PollAnswerService
func NewPollAnswerService(db *gorm.DB) *PollAnswerService {
	return &PollAnswerService{
		repo: repositories.NewPollAnswerRepository(db),
	}
}

// SavePollAnswer сохраняет ответ пользователя на опрос через репозиторий
func (s *PollAnswerService) SavePollAnswer(userID uint, pollID string, optionIDs []int) error {
	return s.repo.SavePollAnswer(userID, pollID, optionIDs)
}
