package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"
	"log"
)

// HandleHelp отправляет сообщение с помощью локализованного текста.
func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}) {
	chatID := message.Chat.ID
	SendMessage(bot, chatID, GetLocalizedString(l, "help_message"))
}

// HandlePollAnswer обрабатывает ответы пользователей на опросы.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - pollAnswer: Ответ пользователя на опрос.
func HandlePollAnswer(bot *tgbotapi.BotAPI, pollAnswer *tgbotapi.PollAnswer, answerService *services.PollAnswerService) {
	log.Printf("Пользователь %d ответил на опрос %s с вариантами %v",
		pollAnswer.User.ID, pollAnswer.PollID, pollAnswer.OptionIDs)

	// Сохранение ответа в базу данных
	err := answerService.SavePollAnswer(uint(pollAnswer.User.ID), pollAnswer.PollID, pollAnswer.OptionIDs)
	if err != nil {
		log.Printf("Ошибка при сохранении ответа на опрос: %v", err)
		return
	}

	log.Println("Ответ на опрос успешно сохранен")
}
