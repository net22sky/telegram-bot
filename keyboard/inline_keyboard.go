package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db"
)

// StartKeyboard создает Inline Keyboard для приветственного сообщения.
func StartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заметки", "notes_menu"),
			tgbotapi.NewInlineKeyboardButtonData("Напоминания", "reminders_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),
		),
	)
}

// NotesKeyboard создает Inline Keyboard для работы с заметками.
func NotesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить заметку", "add_note"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить заметку", "deletes_note"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список заметок", "view_notes"),
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),
		),
	)
}

// RemindersKeyboard создает Inline Keyboard для работы с напоминаниями.
func RemindersKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить напоминание", "add_reminder"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить напоминание", "delete_reminder"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список напоминаний", "view_reminders"),
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),
		),
	)
}

// DeleteNotesKeyboard создает Inline Keyboard для удаления заметок.
func DeleteNotesKeyboard(notes []db.Note) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, note := range notes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d. %s", note.ID, note.Text), fmt.Sprintf("delete_%d", note.ID)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// LanguageKeyboard создает Inline Keyboard для выбора языка.
func LanguageKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "set_language_ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "set_language_en"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
		),
	)
}

// SettingsKeyboard создает Inline Keyboard для настроек пользователя.
func SettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выбрать язык", "choose_language"),
			tgbotapi.NewInlineKeyboardButtonData("Настроить часовой пояс", "choose_timezone"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
		),
	)
}
