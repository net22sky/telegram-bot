package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/keyboard"
	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"
	"log"
	"strings"
)

type Locales map[string]map[string]interface{}

// HandleMessage обрабатывает входящие текстовые сообщения от пользователей.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - update: Входящее обновление от Telegram.
// - locales: Локализованные строки для разных языков.
// - lang: Язык пользователя.
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, locales Locales, lang string, noteService *services.NoteService, userService *services.UserService) {
	l := locales[lang] // Выбираем строки для текущего языка

	if update.Message != nil {
		message := update.Message

		// Проверяем текущее состояние пользователя
		states, exists := state.GetUserState(message.From.ID)

		if exists && states == state.StateAddingNote {
			// Если пользователь добавляет заметку, сохраняем её
			AddNote(bot, message, l, noteService, userService)
			state.SetUserState(int64(message.From.ID), "") // Очищаем состояние
			return
		}

		// Обработка команды /note для создания заметки
		if strings.HasPrefix(message.Text, "/note ") {
			CreateNote(bot, message, l, noteService)
			return
		}

		// Обработка других команд
		switch message.Command() {
		case "notes":
			ViewNotes(bot, message, l, noteService) // Показать список заметок пользователя
		case "start":
			SendStartMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "welcome")) // Отправить приветственное сообщение
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

// CreateNote создает новую заметку для пользователя.
func CreateNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService) {
	parts := strings.SplitN(message.Text, " ", 2)
	if len(parts) < 2 || parts[1] == "" {
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "unknown_command"))
		return
	}

	noteText := parts[1]
	userID := message.From.ID

	// Сохранение заметки в базу данных
	err := noteService.CreateNote(int64(userID), noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "note_creation_error"))
		return
	}

	utils.SendMessage(bot, message.Chat.ID, fmt.Sprintf(utils.GetLocalizedString(l, "note_created"), noteText))
}

// AddNote добавляет заметку для пользователя.
func AddNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService, userService *services.UserService) {
	userID := message.From.ID
	noteText := message.Text

	// Проверяем, существует ли пользователь
	user, err := userService.GetUserByID(int64(userID))
	if user == nil || err != nil {
		// Создаем пользователя, если его нет
		_, err = userService.CreateUser(userID, message.From.UserName, message.From.FirstName)
		if err != nil {
			log.Printf("Ошибка при создании пользователя: %v", err)
			utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "user_creation_error"))
			return
		}
	}

	// Создаем заметку
	err = noteService.CreateNote(int64(userID), noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "note_creation_error"))
		return
	}

	// Уведомляем пользователя об успешном создании заметки
	utils.SendMessage(bot, message.Chat.ID, fmt.Sprintf(utils.GetLocalizedString(l, "note_created"), noteText))
}

// ViewNotes показывает все заметки пользователя.
func ViewNotes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService) {
	userID := message.From.ID

	// Получение списка заметок из базы данных
	notes, err := noteService.GetNotes(userID)
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "no_notes"))
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	utils.SendMessage(bot, message.Chat.ID, utils.GetLocalizedString(l, "notes_list")+response)
}

// SendStartMessage отправляет приветственное сообщение с Inline Keyboard.
func SendStartMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	keyboard := keyboard.StartKeyboard()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
