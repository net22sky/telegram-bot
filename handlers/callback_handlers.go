package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"

	"github.com/net22sky/telegram-bot/db/services"

	"strings"
)

// HandleCallbackQuery обрабатывает нажатия на Inline Keyboard.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - callbackQuery: Входящий callback-запрос.
// - locales: Локализованные строки для разных языков.
// - noteService: Сервис для работы с заметками.
// - userService: Сервис для работы с пользователями.
// - stateManager: Менеджер состояний пользователя.
func HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, locales Locales, noteService *services.NoteService, userService *services.UserService, stateManager *state.StateManager) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	lang := utils.GetUserLanguage(userID, userService)

	l := locales[lang] // Выбираем строки для текущего языка

	switch {
	case strings.HasPrefix(callbackQuery.Data, "delete_"):

		utils.HandleDeleteNote(bot, callbackQuery, l, noteService, userService, chatID, userID)

	case strings.HasPrefix(callbackQuery.Data, "lang_"):
		utils.HandleLanguageChange(bot, callbackQuery, l, userService, chatID, userID)
	case callbackQuery.Data == "cancel":
		// Обработка кнопки "Отмена"
		utils.SendMessage(bot, int64(chatID), utils.GetLocalizedString(l, "action_canceled"))

	default:
		// Обрабатываем остальные действия
		handleDefaultActions(bot, callbackQuery, l, noteService, userService, stateManager, chatID, userID)
	}
}
