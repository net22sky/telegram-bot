package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/keyboard"
	"log"
)

// NotesKeyboard показывает меню заметок.
func NotesKeyboard(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}) {
	chatID := callbackQuery.Message.Chat.ID

	keyboard := keyboard.NotesKeyboard()

	text := GetLocalizedString(l, "notes_menu")

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

	text := GetLocalizedString(l, "notes_menu")

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
		SendMessage(bot, int64(chatID), GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	if len(notes) == 0 {
		SendMessage(bot, int64(chatID), GetLocalizedString(l, "no_notes"))
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. ✍️ %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	SendMessage(bot, int64(chatID), GetLocalizedString(l, "notes_list")+response)
}

// ShowDeleteNotesMenu показывает пользователю меню для удаления заметок.
func ShowDeleteNotesMenu(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}, noteService *services.NoteService) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	// Получаем список заметок пользователя
	notes, err := noteService.GetNotes(int64(userID))
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		SendMessage(bot, chatID, GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	if len(notes) == 0 {
		SendMessage(bot, chatID, GetLocalizedString(l, "no_notes"))
		return
	}

	// Создаем Inline Keyboard для удаления заметок
	keyboard := keyboard.DeleteNotesKeyboard(notes)

	// Отправляем сообщение с клавиатурой
	msg := tgbotapi.NewMessage(chatID, GetLocalizedString(l, "delete_note_prompt"))
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
