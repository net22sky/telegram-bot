package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"

	"log"
)

type Locales map[string]map[string]interface{}

// HandleHelp отправляет сообщение с помощью локализованного текста.
func HandleHelp(bot *tgbotapi.BotAPI, chatID int64, l map[string]interface{}) {
	//chatID := message.Chat.ID
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

func SetLang(bot *tgbotapi.BotAPI, userID int64, chatID int64, langStr string, userService *services.UserService, l map[string]interface{}) {

	err := userService.SetUserLanguage(userID, langStr)

	if err != nil {
		log.Printf("Ошибка при установке языка: %v", err)
		SendMessage(bot, chatID, GetLocalizedString(l, "internal_error")+langStr)
		return
	}

	// Отправляем подтверждение
	SendMessage(bot, chatID, GetLocalizedString(l, "language_set")+langStr)

}
