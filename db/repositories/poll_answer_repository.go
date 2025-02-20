package repositories

import (
	"fmt"
	"github.com/net22sky/telegram-bot/db/models"
	"gorm.io/gorm"
	"log"
)

type PollAnswerRepository struct {
	db *gorm.DB
}

func NewPollAnswerRepository(db *gorm.DB) *PollAnswerRepository {
	return &PollAnswerRepository{db: db}
}

func (r *PollAnswerRepository) SavePollAnswer(userID uint, pollID string, optionIDs []int) error {
	pollAnswer := models.PollAnswer{
		UserID:    userID,
		PollID:    pollID,
		OptionIDs: optionIDs,
	}

	result := r.db.Create(&pollAnswer)
	if result.Error != nil {
		return fmt.Errorf("ошибка при сохранении ответа на опрос: %w", result.Error)
	}

	log.Printf("Ответ на опрос успешно сохранен: %+v", pollAnswer)
	return nil
}
