// keyboard/keyboard.go

package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartKeyboard создает Inline Keyboard для приветственного сообщения.
func StartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить заметку", "note"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить заметку", "dellnote"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список заметок", "notes"),
		),
	)
}

// DeleteNotesKeyboard создает Inline Keyboard для удаления заметок.
func DeleteNotesKeyboard(notes []Note) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Создаем кнопки для каждой заметки
	for _, note := range notes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d. %s", note.ID, note.Text), fmt.Sprintf("delete_%d", note.ID)),
		))
	}

	// Добавляем кнопку "Отмена"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Note представляет заметку пользователя.
type Note struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
}
