package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/keyboard"
	"log"
)

// SendMessage отправляет текстовое сообщение в указанный чат.
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// SendPoll отправляет опрос в указанный чат.
func SendPoll(bot *tgbotapi.BotAPI, chatID int64) {
	poll := tgbotapi.NewPoll(chatID, "Какой ваш любимый язык программирования?", "Go", "Python", "JavaScript", "Java")
	poll.IsAnonymous = false
	if _, err := bot.Send(poll); err != nil {
		log.Printf("Ошибка при отправке опроса: %v", err)
	}
}

// SendStartMessage отправляет приветственное сообщение с Inline Keyboard.
func SendStartMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	keyboard := keyboard.StartKeyboard()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
