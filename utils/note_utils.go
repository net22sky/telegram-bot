package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"
	"log"
	"strconv"
	"strings"
)

// DeleteNote удаляет заметку по её ID, если она принадлежит указанному пользователю.
func DeleteNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService) {
	parts := strings.SplitN(message.Text, " ", 2)
	if len(parts) < 2 || parts[1] == "" {
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "unknown_command"))
		return
	}

	noteIDStr := parts[1]
	userID := message.From.ID
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil || noteID <= 0 {
		log.Printf("Ошибка при преобразовании ID заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "invalid_note_id"))
		return
	}

	err = noteService.DeleteNoteByID(uint(noteID), int64(userID))
	if err != nil {
		log.Printf("Ошибка при удалении заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_deletion_error"))
		return
	}

	SendMessage(bot, message.Chat.ID, fmt.Sprintf(GetLocalizedString(l, "note_deleted"), noteID))
}

// ViewNotes показывает все заметки пользователя.
func ViewNotes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService) {
	userID := message.From.ID

	// Получение списка заметок из базы данных
	notes, err := noteService.GetNotes(userID)
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	if len(notes) == 0 {
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "no_notes"))
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "notes_list")+response)
}

// CreateNote создает новую заметку для пользователя.
func CreateNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService) {
	parts := strings.SplitN(message.Text, " ", 2)
	if len(parts) < 2 || parts[1] == "" {
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "unknown_command"))
		return
	}

	noteText := parts[1]
	userID := message.From.ID

	// Сохранение заметки в базу данных
	err := noteService.CreateNote(int64(userID), noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_creation_error"))
		return
	}

	SendMessage(bot, message.Chat.ID, fmt.Sprintf(GetLocalizedString(l, "note_created"), noteText))
}

// AddNote добавляет заметку для пользователя.
func AddNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService) {
	userID := message.From.ID
	noteText := message.Text

	// Проверяем, существует ли пользователь
	user, err := userService.GetUserByID(int64(userID))
	if user == nil || err != nil {
		// Создаем пользователя, если его нет
		_, err = userService.CreateUser(userID, message.From.UserName, message.From.FirstName)
		if err != nil {
			log.Printf("Ошибка при создании пользователя: %v", err)
			SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "user_creation_error"))
			return
		}
	}

	// Создаем заметку
	err = noteService.CreateNote(int64(userID), noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_creation_error"))
		return
	}

	// Уведомляем пользователя об успешном создании заметки
	SendMessage(bot, message.Chat.ID, fmt.Sprintf(GetLocalizedString(l, "note_created"), noteText))
}
