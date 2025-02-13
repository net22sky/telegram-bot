// utils/utils.go
package utils

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

func SendPoll(bot *tgbotapi.BotAPI, chatID int64) {
	poll := tgbotapi.NewPoll(chatID, "Какой ваш любимый язык программирования?", "Go", "Python", "JavaScript", "Java")
	poll.IsAnonymous = false // Опрос не анонимный
	if _, err := bot.Send(poll); err != nil {
		log.Printf("Ошибка при отправке опроса: %v", err)
	}
}
