package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/state"
	"log"
	"strconv"
	"strings"
)

// handleUserState обрабатывает текущее состояние пользователя.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - l: Локализованные строки для текущего языка.
// - noteService: Сервис для работы с заметками.
// - userService: Сервис для работы с пользователями.
// - stateManager: Менеджер состояний пользователя.
// Возвращает:
// - true, если состояние было обработано; false, если состояние отсутствует или не требует обработки.
func HandleUserState(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService, stateManager *state.StateManager) bool {
	states, exists := stateManager.GetUserState(message.From.ID)
	if exists && states == state.StateAddingNote {
		// Если пользователь находится в состоянии добавления заметки, сохраняем её
		AddNote(bot, message, l, noteService, userService)
		stateManager.DeleteUserState(int64(message.From.ID)) // Очищаем состояние
		return true
	}
	return false
}

// handleDeleteNote обрабатывает удаление заметки.
func HandleDeleteNote(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService, chatID int64, userID int64) {
	noteIDStr := strings.TrimPrefix(callbackQuery.Data, "delete_")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil || noteID <= 0 {
		log.Printf("Ошибка при парсинге ID заметки: %v", err)
		SendMessage(bot, chatID, GetLocalizedString(l, "invalid_note_id"))
		return
	}

	note, err := noteService.GetNoteByID(int64(noteID))
	if err != nil {
		log.Printf("Ошибка при получении заметки: %v", err)
		SendMessage(bot, chatID, GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	user, err := userService.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		return
	}

	if note == nil || note.UserID != uint(user.ID) {
		SendMessage(bot, chatID, GetLocalizedString(l, "note_not_found"))
		return
	}

	err = noteService.DeleteNoteByID(uint(noteID), int64(userID))
	if err != nil {
		log.Printf("Ошибка при удалении заметки: %v", err)
		SendMessage(bot, chatID, GetLocalizedString(l, "note_deletion_error"))
		return
	}

	SendMessage(bot, chatID, fmt.Sprintf(GetLocalizedString(l, "note_deleted"), noteID))
}

// handleLanguageChange обрабатывает изменение языка.
func HandleLanguageChange(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}, userService *services.UserService, chatID int64, userID int64) {
	langStr := strings.TrimPrefix(callbackQuery.Data, "lang_")
	if len(langStr) > 0 {
		SetLang(bot, int64(userID), chatID, langStr, userService, l)
	}
}
