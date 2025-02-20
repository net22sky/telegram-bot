package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"

	"github.com/net22sky/telegram-bot/db/services"
	"log"
	"strconv"
	"strings"
)

// HandleCallbackQuery обрабатывает нажатия на Inline Keyboard.
func HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, locales Locales, lang string, noteService *services.NoteService, userService *services.UserService, stateManager  *state.StateManager ) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	l := locales[lang] // Выбираем строки для текущего языка

	switch {
	case strings.HasPrefix(callbackQuery.Data, "delete_"):
		// Извлекаем ID заметки из данных кнопки
		noteIDStr := strings.TrimPrefix(callbackQuery.Data, "delete_")
		noteID, err := strconv.Atoi(noteIDStr)
		if err != nil || noteID <= 0 {
			log.Printf("Ошибка при парсинге ID заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "invalid_note_id"))
			return
		}

		// Получаем заметку для проверки владельца
		note, err := noteService.GetNoteByID(int64(noteID))
		if err != nil {
			log.Printf("Ошибка при получении заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "note_retrieval_error"))
			return
		}

		// Получаем пользователя по Telegram ID
		user, err := userService.GetUserByID(userID)
		if err != nil {
			log.Printf("ошибка при получении пользователя: %v", err)
			return
		}

		if note == nil || note.UserID != uint(user.ID) {
			utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "note_not_found"))
			return
		}

		// Удаляем заметку
		err = noteService.DeleteNoteByID(uint(noteID), int64(userID))
		if err != nil {
			log.Printf("Ошибка при удалении заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "note_deletion_error"))
			return
		}

		// Отправляем сообщение об успешном удалении
		utils.SendMessage(bot, int64(chatID), fmt.Sprintf(utils.GetLocalizedString(l, "note_deleted"), noteID))

	case callbackQuery.Data == "cancel":
		// Обработка кнопки "Отмена"
		utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "action_canceled"))

	default:
		// Обрабатываем остальные действия
		switch callbackQuery.Data {
		case "notes_menu":
			utils.NotesKeyboard(bot, callbackQuery, l)
		case "reminders_menu":
			utils.RemindersKeyboard(bot, callbackQuery, l)
		case "add_note":
			// Переходим в режим добавления заметки
			
			stateManager.SetUserState(userID, state.StateAddingNote)
			utils.SendMessage(bot, chatID, utils.GetLocalizedString(l, "enter_note_text"))

		case "view_notes":
			utils.ViewNotesKeyboard(bot, callbackQuery, l, noteService)
		case "deletes_note":
			utils.ShowDeleteNotesMenu(bot, callbackQuery, l, noteService)
		default:
			utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "unknown_action"))
		}
	}
}
