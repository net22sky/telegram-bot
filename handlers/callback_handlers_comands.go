package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"
)

func handleDefaultActions(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService, stateManager *state.StateManager, chatID int64, userID int64) {
	switch callbackQuery.Data {
	case "main_menu":
		utils.SendStartMessage(bot, chatID, utils.GetLocalizedString(l, "welcome"))
	case "help":
		utils.HandleHelp(bot, chatID, l)
	case "notes_menu":
		utils.NotesKeyboard(bot, callbackQuery, l)
	case "reminders_menu":
		utils.RemindersKeyboard(bot, callbackQuery, l)
	case "add_note":
		stateManager.SetUserState(userID, state.StateAddingNote)
		utils.SendMessage(bot, chatID, utils.GetLocalizedString(l, "enter_note_text"))
	case "view_notes":
		utils.ViewNotesKeyboard(bot, callbackQuery, l, noteService)
	case "deletes_note":
		utils.ShowDeleteNotesMenu(bot, callbackQuery, l, noteService)
	case "add_reminder":
		utils.AddReminderKeyboard(bot, chatID, l)
	case "category_subscription":
		utils.SendMessage(bot, chatID, "Вы выбрали: Подписка")
	case "category_birthday":
		utils.SendMessage(bot, chatID, "Вы выбрали: День рождения")
	case "category_loans":
		utils.SendMessage(bot, chatID, "Вы выбрали: Кредиты")
	case "category_utilities":
		utils.SendMessage(bot, chatID, "Вы выбрали: ЖКХ")
	default:
		utils.SendMessage(bot, chatID, utils.GetLocalizedString(l, "unknown_action"))
	}
}
