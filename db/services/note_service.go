package services

import (
	"fmt"
	"github.com/net22sky/telegram-bot/db/models"
	"github.com/net22sky/telegram-bot/db/repositories"
)

type NoteService struct {
	noteRepo *repositories.NoteRepository
	userRepo *repositories.UserRepository
}

func NewNoteService(noteRepo *repositories.NoteRepository, userRepo *repositories.UserRepository) *NoteService {
	return &NoteService{
		noteRepo: noteRepo,
		userRepo: userRepo}
}

func (s *NoteService) CreateNote(telegramID int64, text string) error {
	user, err := s.userRepo.GetUserByID(telegramID)
	if err != nil {
		return err
	}

	return s.noteRepo.CreateNote(user.ID, text)
}

func (s *NoteService) GetNotes(userID int64) ([]models.Note, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.noteRepo.GetNotes(user.ID)
}

func (s *NoteService) DeleteNoteByID(noteID uint, userID int64) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.noteRepo.DeleteNoteByID(noteID, user.ID)
}

func (s *NoteService) GetNoteByID(noteID int64) (*models.Note, error) {

	note, err := s.noteRepo.GetNoteByID(noteID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении заметки: %w", err)
	}
	if note == nil {
		return nil, fmt.Errorf("заметка не найдена")
	}
	return note, nil
}
