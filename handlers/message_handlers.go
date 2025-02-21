package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"

	"strings"
)

type Locales map[string]map[string]interface{}

// HandleMessage обрабатывает входящие текстовые сообщения от пользователей.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - update: Входящее обновление от Telegram.
// - locales: Локализованные строки для разных языков.
// - noteService: Сервис для работы с заметками.
// - userService: Сервис для работы с пользователями.
// - stateManager: Менеджер состояний пользователя.
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, locales Locales, noteService *services.NoteService, userService *services.UserService, stateManager *state.StateManager) {

	//chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	lang := utils.GetUserLanguage(userID, userService)

	l := locales[lang] // Выбираем строки для текущего языка

	if update.Message != nil {
		message := update.Message

		// Проверяем текущее состояние пользователя
		// Обработка состояний пользователя (например, добавление заметки)
		if utils.HandleUserState(bot, message, l, noteService, userService, stateManager) {
			return
		}

		// Обработка команды /note для создания заметки
		if strings.HasPrefix(message.Text, "/note ") {
			utils.CreateNote(bot, message, l, noteService)
			return
		}

		// Обработка команд пользователя
		if handleCommands(bot, message, l, noteService, userService) {
			return
		}
		// Если команда не распознана, отправляем сообщение о неизвестной команде
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "unknown_command"))

	}
}
