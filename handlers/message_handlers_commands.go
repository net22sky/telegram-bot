package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/utils"

	"strings"
)

func handleCommands(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService) bool {
	// Обработка команды /note
	if strings.HasPrefix(message.Text, "/note ") {
		utils.CreateNote(bot, message, l, noteService)
		return true
	}
	// Обработка других команд
	switch message.Command() {
	case "notes":
		utils.ViewNotes(bot, message, l, noteService) // Показать список заметок пользователя
	case "start":
		utils.SendStartMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "welcome")) // Отправить приветственное сообщение
	case "help":
		utils.HandleHelp(bot, message.Chat.ID, l) // Обработка команды /help
	case "poll":
		utils.SendPoll(bot, message.Chat.ID) // Создать опрос
	case "dellnote":
		utils.DeleteNote(bot, message, l, noteService) // Удалить заметку
	case "settings":
		utils.ViewSendSettingsKeyboard(bot, message.Chat.ID, l) // Отправить настройки
	case "clear":
		// Очистка чата
		utils.ClearChat(bot, message.Chat.ID, l)

	default:
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "unknown_command")) // Сообщение о неизвестной команде
	}
	return true
}
