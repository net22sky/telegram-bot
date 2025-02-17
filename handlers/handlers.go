package handlers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db"
	"github.com/net22sky/telegram-bot/keyboard"
	"github.com/net22sky/telegram-bot/state"
	"github.com/net22sky/telegram-bot/utils"
	"log"
	"strconv"
	"strings"
)

// Locales содержит строки для разных языков, используемые для локализации сообщений бота.
type Locales map[string]map[string]string

// HandleMessage обрабатывает входящие текстовые сообщения от пользователей.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - locales: Локализованные строки для разных языков.
// - lang: Язык пользователя.
// func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, locales Locales, lang string) {
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, locales Locales, lang string) {

	//log.Printf("[%s] %s", message.From.UserName, message.Text)

	l := locales[lang] // Выбираем строки для текущего языка

	// Проверяем, является ли обновление CallbackQuery
	if update.CallbackQuery != nil {
		HandleCallbackQuery(bot, update.CallbackQuery, locales, lang)

	}

	// Если это обычное сообщение
	if update.Message != nil {
		message := update.Message

		// Проверяем текущее состояние пользователя
		states, exists := state.GetUserState(message.From.ID)

		if exists && states == state.StateAddingNote {
			// Если пользователь добавляет заметку, сохраняем её
			AddNote(bot, message, l)
			state.SetUserState(int64(message.From.ID), "") // Очищаем состояние
			return
		}

		// Обработка команды /note для создания заметки
		if strings.HasPrefix(message.Text, "/note ") {
			CreateNote(bot, message, l)
			return
		}

		// Обработка других команд
		switch message.Command() {
		case "notes":
			ViewNotes(bot, message, l) // Показать список заметок пользователя
		case "start":
			SendStartMessage(bot, message.Chat.ID, l["welcome"]) // Отправить приветственное сообщение
		case "help":
			utils.HandleHelp(bot, message, l) // Обработка команды /help
		case "poll":
			utils.SendPoll(bot, message.Chat.ID) // Создать опрос
		case "dellnote":
			utils.DeleteNote(bot, message, l) // Создать опрос
		default:
			utils.SendMessage(bot, message.Chat.ID, l["unknown_command"]) // Сообщение о неизвестной команде
		}
	}

}

// CreateNote создает новую заметку для пользователя.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - l: Локализованные строки для текущего языка.
func CreateNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string) {
	parts := strings.SplitN(message.Text, " ", 2)
	if len(parts) < 2 || parts[1] == "" {
		utils.SendMessage(bot, message.Chat.ID, l["unknown_command"])
		return
	}

	noteText := parts[1]
	userID := message.From.ID

	// Сохранение заметки в базу данных
	err := db.CreateNote(int64(userID), noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		utils.SendMessage(bot, message.Chat.ID, l["note_creation_error"])
		return
	}

	utils.SendMessage(bot, message.Chat.ID, fmt.Sprintf(l["note_created"], noteText))
}

// HandlePollAnswer обрабатывает ответы пользователей на опросы.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - pollAnswer: Ответ пользователя на опрос.
func HandlePollAnswer(bot *tgbotapi.BotAPI, pollAnswer *tgbotapi.PollAnswer) {
	log.Printf("Пользователь %d ответил на опрос %s с вариантами %v",
		pollAnswer.User.ID, pollAnswer.PollID, pollAnswer.OptionIDs)

	// Сохранение ответа в базу данных
	err := db.SavePollAnswer(uint(pollAnswer.User.ID), pollAnswer.PollID, pollAnswer.OptionIDs)
	if err != nil {
		log.Printf("Ошибка при сохранении ответа на опрос: %v", err)
		return
	}

	log.Println("Ответ на опрос успешно сохранен")
}

// HandleCallbackQuery обрабатывает нажатия на Inline Keyboard.
func HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, locales Locales, lang string) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	l := locales[lang] // Выбираем строки для текущего языка

	//utils.SendMessage(bot, int64(chatID), "HandleCallbackQuery")
	// Обрабатываем действие по нажатию кнопки
	switch {
	case strings.HasPrefix(callbackQuery.Data, "delete_"):
		// Извлекаем ID заметки из данных кнопки
		noteIDStr := strings.TrimPrefix(callbackQuery.Data, "delete_")
		noteID, err := strconv.Atoi(noteIDStr)
		if err != nil || noteID <= 0 {
			log.Printf("Ошибка при парсинге ID заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), l["invalid_note_id"])
			return
		}

		// Получаем заметку для проверки владельца
		note, err := db.GetNoteByID(int64(noteID))
		if err != nil {
			log.Printf("Ошибка при получении заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), l["note_retrieval_error"])
			return
		}

		// Получаем пользователя по Telegram ID
		user, err := db.GetUserByID(userID)
		if err != nil {
			log.Printf("ошибка при получении пользователя: %v", err)
			return
		}

		if note == nil || note.UserID != uint(user.ID) {
			utils.SendMessage(bot, int64(chatID), l["note_not_found"])
			return
		}

		// Удаляем заметку
		err = db.DeleteNoteByID(uint(noteID), int64(userID))
		if err != nil {
			log.Printf("Ошибка при удалении заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), l["note_deletion_error"])
			return
		}

		// Отправляем сообщение об успешном удалении
		utils.SendMessage(bot, int64(chatID), fmt.Sprintf(l["note_deleted"], noteID))

	case callbackQuery.Data == "cancel":
		// Обработка кнопки "Отмена"
		utils.SendMessage(bot, int64(chatID), l["action_canceled"])

	default:
		// Обрабатываем остальные действия
		switch callbackQuery.Data {
		case "add_note":
			// Переходим в режим добавления заметки
			state.SetUserState(userID, state.StateAddingNote)
			utils.SendMessage(bot, chatID, l["enter_note_text"])

		case "view_notes":
			ViewNotesKeyboard(bot, callbackQuery, l)
		case "deletes_note":
			ShowDeleteNotesMenu(bot, callbackQuery, l)
		default:
			utils.SendMessage(bot, int64(chatID), l["unknown_action"])
		}
	}
}

// SendStartMessage отправляет приветственное сообщение с Inline Keyboard.
func SendStartMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	// Создаем Inline Keyboard через пакет keyboard
	keyboard := keyboard.StartKeyboard()

	// Отправляем приветственное сообщение с клавиатурой
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// AddNote добавляет заметку для пользователя.
func AddNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string) {
	userID := message.From.ID
	noteText := message.Text

	// Проверяем, существует ли пользователь
	user, err := db.GetUserByID(int64(userID))
	if user == nil || err != nil {
		// Создаем пользователя, если его нет
		_, err = db.CreateUser(userID, message.From.UserName, message.From.FirstName)
		if err != nil {
			log.Printf("Ошибка при создании пользователя: %v", err)
			utils.SendMessage(bot, message.Chat.ID, l["user_creation_error"])
			return
		}
	}

	// Создаем заметку
	err = db.CreateNote(int64(userID), noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		utils.SendMessage(bot, message.Chat.ID, l["note_creation_error"])
		return
	}

	// Уведомляем пользователя об успешном создании заметки
	utils.SendMessage(bot, message.Chat.ID, fmt.Sprintf(l["note_created"], noteText))
}

// ViewNotesKeyboard показывает все заметки пользователя.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - l: Локализованные строки для текущего языка.
func ViewNotesKeyboard(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]string) {
	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	// Получение списка заметок из базы данных
	notes, err := db.GetNotes(int64(userID))
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, int64(chatID), l["note_retrieval_error"])
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, int64(chatID), l["no_notes"])
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. ✍️ %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	utils.SendMessage(bot, int64(chatID), l["notes_list"]+response)
}

// cancelAction отменяет текущее действие пользователя.
func CancelAction(bot *tgbotapi.BotAPI, chatID int64, userID int64, l map[string]string) {

	state.DeleteUserState(userID)
	utils.SendMessage(bot, chatID, l["action_canceled"])
}

// ViewNotes показывает все заметки пользователя.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - l: Локализованные строки для текущего языка.
func ViewNotes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string) {
	userID := message.From.ID

	// Получение списка заметок из базы данных
	notes, err := db.GetNotes(userID)
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, message.Chat.ID, l["note_retrieval_error"])
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, message.Chat.ID, l["no_notes"])
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	utils.SendMessage(bot, message.Chat.ID, l["notes_list"]+response)
}

// ShowDeleteNotesMenu показывает пользователю меню для удаления заметок.
func ShowDeleteNotesMenu(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, l map[string]string) {
	////userID := callbackQuery.From.ID
	//chatID := callbackQuery.Chat.ID

	chatID := callbackQuery.Message.Chat.ID
	userID := callbackQuery.From.ID

	// Получаем список заметок пользователя
	notes, err := db.GetNotes(int64(userID))
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, chatID, l["note_retrieval_error"])
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, chatID, l["no_notes"])
		return
	}

	// Создаем Inline Keyboard для удаления заметок
	keyboard := keyboard.DeleteNotesKeyboard(notes)

	// Отправляем сообщение с клавиатурой
	msg := tgbotapi.NewMessage(chatID, l["delete_note_prompt"])
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
