package repositories

import (
	"fmt"
	"github.com/net22sky/telegram-bot/db/models"
	"gorm.io/gorm"
	"log"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) CreateNote(userID uint, text string) error {
	note := models.Note{
		UserID: userID,
		Text:   text,
	}

	result := r.db.Create(&note)
	if result.Error != nil {
		return fmt.Errorf("ошибка при создании заметки: %w", result.Error)
	}

	log.Printf("Заметка успешно создана: %+v", note)
	return nil
}

func (r *NoteRepository) GetNotes(userID uint) ([]models.Note, error) {
	var notes []models.Note

	result := r.db.Where("user_id = ?", userID).Find(&notes)
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при получении заметок: %w", result.Error)
	}

	if len(notes) == 0 {
		log.Printf("У пользователя %d нет заметок", userID)
	}

	return notes, nil
}
func (r *NoteRepository) GetNoteByID(noteID int64) (*models.Note, error) {
	var note models.Note

	// Находим заметку по ID
	result := r.db.First(&note, int64(noteID))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Заметка с ID %d не найдена", noteID)
			return nil, nil // Заметка не найдена
		}
		return nil, fmt.Errorf("ошибка при получении заметки: %w", result.Error)
	}

	log.Printf("Заметка с ID %d успешно загружена: %+v", noteID, note)
	return &note, nil
}

func (r *NoteRepository) DeleteNoteByID(noteID, userID uint) error {
	var note models.Note

	result := r.db.First(&note, noteID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("заметка с ID %d не найдена", noteID)
		}
		return fmt.Errorf("ошибка при получении заметки: %w", result.Error)
	}

	if note.UserID != userID {
		return fmt.Errorf("заметка с ID %d не принадлежит пользователю %d", noteID, userID)
	}

	result = r.db.Delete(&note)
	if result.Error != nil {
		return fmt.Errorf("ошибка при удалении заметки: %w", result.Error)
	}

	log.Printf("Заметка с ID %d успешно удалена", noteID)
	return nil
}
