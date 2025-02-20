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
// - lang: Язык пользователя.
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, locales Locales, lang string, noteService *services.NoteService, userService *services.UserService,stateManager  *state.StateManager) {
	l := locales[lang] // Выбираем строки для текущего языка

	if update.Message != nil {
		message := update.Message

		
		// Проверяем текущее состояние пользователя
		states, exists := stateManager.GetUserState(message.From.ID)

		if exists && states == state.StateAddingNote {
			// Если пользователь добавляет заметку, сохраняем её
			utils.AddNote(bot, message, l, noteService, userService)
			stateManager.DeleteUserState(int64(message.From.ID)) // Очищаем состояние
			return
		}

		// Обработка команды /note для создания заметки
		if strings.HasPrefix(message.Text, "/note ") {
			utils.CreateNote(bot, message, l, noteService)
			return
		}

		// Обработка других команд
		switch message.Command() {
		case "notes":
			utils.ViewNotes(bot, message, l, noteService) // Показать список заметок пользователя
		case "start":
			utils.SendStartMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "welcome")) // Отправить приветственное сообщение
		case "help":
			utils.HandleHelp(bot, message, l) // Обработка команды /help
		case "poll":
			utils.SendPoll(bot, message.Chat.ID) // Создать опрос
		case "dellnote":
			utils.DeleteNote(bot, message, l, noteService) // Удалить заметку

		default:
			utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "unknown_command")) // Сообщение о неизвестной команде
		}
	}
}

// handleUserState обрабатывает состояние пользователя (например, добавление заметки).
func handleUserState(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService,stateManager  *state.StateManager) bool {

	states, exists := stateManager.GetUserState(message.From.ID)
	if exists && states == state.StateAddingNote {
		utils.AddNote(bot, message, l, noteService, userService)
		stateManager.DeleteUserState(int64(message.From.ID))
		return true
	}
	return false
}
