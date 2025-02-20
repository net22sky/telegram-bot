package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/keyboard"
	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"

	"github.com/net22sky/telegram-bot/db/services"
	"log"
	"strconv"
	"strings"
)

// HandleCallbackQuery обрабатывает нажатия на Inline Keyboard.
func HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, locales Locales, lang string, noteService *services.NoteService, userService *services.UserService) {
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
			NotesKeyboard(bot, callbackQuery, l)
		case "reminders_menu":
			RemindersKeyboard(bot, callbackQuery, l)
		case "add_note":
			// Переходим в режим добавления заметки
			state.SetUserState(userID, state.StateAddingNote)
			utils.SendMessage(bot, chatID, utils.GetLocalizedString(l, "enter_note_text"))

		case "view_notes":
			ViewNotesKeyboard(bot, callbackQuery, l, noteService)
		case "deletes_note":
			ShowDeleteNotesMenu(bot, callbackQuery, l, noteService)
		default:
			utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "unknown_action"))
		}
	}
}

// NotesKeyboard показывает меню заметок.
func NotesKeyboard(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}) {
	chatID := callbackQuery.Message.Chat.ID

	keyboard := keyboard.NotesKeyboard()

	text := utils.GetLocalizedString(l, "notes_menu")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// RemindersKeyboard показывает меню напоминаний.
func RemindersKeyboard(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}) {
	chatID := callbackQuery.Message.Chat.ID

	keyboard := keyboard.RemindersKeyboard()

	text := utils.GetLocalizedString(l, "notes_menu")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// ViewNotesKeyboard показывает все заметки пользователя.
func ViewNotesKeyboard(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}, noteService *services.NoteService) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	// Получение списка заметок из базы данных
	notes, err := noteService.GetNotes(int64(userID))
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "no_notes"))
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. ✍️ %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "notes_list")+response)
}

// ShowDeleteNotesMenu показывает пользователю меню для удаления заметок.
func ShowDeleteNotesMenu(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}, noteService *services.NoteService) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	// Получаем список заметок пользователя
	notes, err := noteService.GetNotes(int64(userID))
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, chatID, utils.GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, chatID, utils.GetLocalizedString(l, "no_notes"))
		return
	}

	// Создаем Inline Keyboard для удаления заметок
	keyboard := keyboard.DeleteNotesKeyboard(notes)

	// Отправляем сообщение с клавиатурой
	msg := tgbotapi.NewMessage(chatID, utils.GetLocalizedString(l, "delete_note_prompt"))
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
