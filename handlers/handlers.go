package handlers

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/keyboard"
	"github.com/net22sky/telegram-bot/mysql"
	"github.com/net22sky/telegram-bot/utils"
	"log"
	"strconv"
	"strings"
)

// Locales содержит строки для разных языков, используемые для локализации сообщений бота.
type Locales map[string]map[string]string

var userStates = make(map[int64]string) // Хранилище состояний пользователей
const (
	StateIdle       = "idle"        // Пользователь находится в режиме ожидания
	StateAddingNote = "adding_note" // Пользователь добавляет заметку
)

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
		return
	}

	// Если это обычное сообщение
	if update.Message != nil {
		message := update.Message
		// Обработка команды /note для создания заметки
		if strings.HasPrefix(message.Text, "/note ") {
			CreateNote(bot, message, l)
			return
		}

		// Проверяем текущее состояние пользователя
		state, exists := userStates[message.From.ID]
		if exists && state == StateAddingNote {
			// Если пользователь добавляет заметку, сохраняем её
			AddNote(bot, message, l)
			delete(userStates, message.From.ID) // Очищаем состояние
			return
		}

		// Обработка других команд
		switch message.Command() {
		case "notes":
			ViewNotes(bot, message, l) // Показать список заметок пользователя
		case "start":
			SendStartMessage(bot, message.Chat.ID, l["welcome"]) // Отправить приветственное сообщение
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
	err := mysql.CreateNote(userID, noteText)
	if err != nil {
		log.Printf("Ошибка при создании заметки: %v", err)
		utils.SendMessage(bot, message.Chat.ID, l["note_creation_error"])
		return
	}

	utils.SendMessage(bot, message.Chat.ID, fmt.Sprintf(l["note_created"], noteText))
}

// ViewNotes показывает все заметки пользователя.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - l: Локализованные строки для текущего языка.
func ViewNotes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string) {
	userID := message.From.ID

	// Получение списка заметок из базы данных
	notes, err := mysql.GetNotes(userID)
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

// HandlePollAnswer обрабатывает ответы пользователей на опросы.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - pollAnswer: Ответ пользователя на опрос.
func HandlePollAnswer(bot *tgbotapi.BotAPI, pollAnswer *tgbotapi.PollAnswer) {
	log.Printf("Пользователь %d ответил на опрос %s с вариантами %v",
		pollAnswer.User.ID, pollAnswer.PollID, pollAnswer.OptionIDs)

	// Сохранение ответа в базу данных
	err := mysql.SavePollAnswer(pollAnswer.User.ID, pollAnswer.PollID, pollAnswer.OptionIDs)
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

	utils.SendMessage(bot, int64(chatID), "invalid_note_id")
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
		note, err := mysql.GetNoteByID(int64(noteID))
		if err != nil {
			log.Printf("Ошибка при получении заметки: %v", err)
			utils.SendMessage(bot, int64(chatID), l["note_retrieval_error"])
			return
		}

		if note == nil || note.UserID != userID {
			utils.SendMessage(bot, int64(chatID), l["note_not_found"])
			return
		}

		// Удаляем заметку
		err = mysql.DeleteNoteByID(noteID)
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
			utils.SendMessage(bot, int64(chatID), l["add_note_prompt"])
		case "view_notes":
			ViewNotesKeyboard(bot, int64(userID),int64(chatID), l)
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

	// Создаем заметку
	err := mysql.CreateNote(userID, noteText)
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
func ViewNotesKeyboard(bot *tgbotapi.BotAPI, userID int64, ChatID int64, l map[string]string) {


	// Получение списка заметок из базы данных
	notes, err := mysql.GetNotes(userID)
	if err != nil {
		log.Printf("Ошибка при получении заметок: %v", err)
		utils.SendMessage(bot, ChatID, l["note_retrieval_error"])
		return
	}

	if len(notes) == 0 {
		utils.SendMessage(bot, ChatID, l["no_notes"])
		return
	}

	// Формирование списка заметок для отправки пользователю
	var response string
	for i, note := range notes {
		response += fmt.Sprintf("%d. %s (ID: %d)\n", i+1, note.Text, note.ID)
	}

	utils.SendMessage(bot, ChatID, l["notes_list"]+response)
}
