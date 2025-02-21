// ClearChat очищает чат от сообщений.
package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"log"
)

func ClearChat(bot *tgbotapi.BotAPI, chatID int64, l map[string]interface{}) {
	// Проверяем, является ли чат группой или супергруппой
	chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}})
	if err != nil {
		log.Printf("Ошибка при получении информации о чате: %v", err)
		SendMessage(bot, chatID, GetLocalizedString(l, "clear_chat_error"))
		return
	}

	// Очистка возможна только в группах и супергруппах
	if chat.Type != "group" && chat.Type != "supergroup" {
		SendMessage(bot, chatID, GetLocalizedString(l, "clear_chat_not_supported"))
		return
	}

	// Получаем последние 100 сообщений
	deleteMessages := []int{}
	for i := 0; i < 100; i++ {
		deleteMessages = append(deleteMessages, i)
	}

	// Удаляем сообщения
	for _, msgID := range deleteMessages {
		config := tgbotapi.NewDeleteMessage(chatID, msgID)
		if _, err := bot.Request(config); err != nil {
			log.Printf("Ошибка при удалении сообщения %d: %v", msgID, err)
		}
	}

	// Отправляем подтверждение
	SendMessage(bot, chatID, GetLocalizedString(l, "clear_chat_success"))
}
