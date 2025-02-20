package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// HandleHelp отправляет сообщение с помощью локализованного текста.
func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}) {
	chatID := message.Chat.ID
	SendMessage(bot, chatID, GetLocalizedString(l, "help_message"))
}
