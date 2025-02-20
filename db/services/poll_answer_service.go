package services

import (
	"github.com/net22sky/telegram-bot/db/repositories"
)

// PollAnswerService предоставляет методы для работы с ответами на опросы
type PollAnswerService struct {
	pollAnswerRepo *repositories.PollAnswerRepository
}

// NewPollAnswerService создает новый экземпляр PollAnswerService
func NewPollAnswerService(pollAnswerRepo *repositories.PollAnswerRepository) *PollAnswerService {
	return &PollAnswerService{
		pollAnswerRepo: pollAnswerRepo,
	}
}

// SavePollAnswer сохраняет ответ пользователя на опрос через репозиторий
func (s *PollAnswerService) SavePollAnswer(userID uint, pollID string, optionIDs []int) error {
	return s.pollAnswerRepo.SavePollAnswer(userID, pollID, optionIDs)
}
