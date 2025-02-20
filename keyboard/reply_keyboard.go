package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// CreateNumberKeyboard создает Reply Keyboard с числовыми кнопками.
func CreateNumberKeyboard(start, end, rowSize int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for i := start; i <= end; i += rowSize {
		var row []tgbotapi.KeyboardButton

		for j := i; j < i+rowSize && j <= end; j++ {
			row = append(row, tgbotapi.NewKeyboardButton(fmt.Sprintf("%d", j)))
		}

		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

// CreateMonthKeyboard создает Reply Keyboard с кнопками месяцев.
func CreateMonthKeyboard(months []string) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for i := 0; i < len(months); i += 4 {
		end := i + 4
		if end > len(months) {
			end = len(months)
		}

		row := make([]tgbotapi.KeyboardButton, 0)
		for _, month := range months[i:end] {
			row = append(row, tgbotapi.NewKeyboardButton(month))
		}

		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

// RemoveReplyKeyboard удаляет текущую Reply Keyboard.
func RemoveReplyKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при удалении клавиатуры: %v", err)
	}
}

// ConfirmationKeyboard создает Reply Keyboard для подтверждения действий.
func ConfirmationKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сохранить"),
			tgbotapi.NewKeyboardButton("Отменить"),
		),
	)
	keyboard.ResizeKeyboard = true

	return keyboard
}

// ReminderTypesKeyboard создает Reply Keyboard для выбора типа напоминания.
func ReminderTypesKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("День рождения"),
			tgbotapi.NewKeyboardButton("Credit"),
			tgbotapi.NewKeyboardButton("Подписка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отменить"),
		),
	)
	keyboard.ResizeKeyboard = true

	return keyboard
}
