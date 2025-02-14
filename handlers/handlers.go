package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/mysql"
	"github.com/net22sky/telegram-bot/utils"
)

// Locales содержит строки для разных языков, используемые для локализации сообщений бота.
type Locales map[string]map[string]string

// HandleMessage обрабатывает входящие текстовые сообщения от пользователей.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - message: Входящее сообщение от пользователя.
// - locales: Локализованные строки для разных языков.
// - lang: Язык пользователя.
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, locales Locales, lang string) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	l := locales[lang] // Выбираем строки для текущего языка

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
		utils.SendMessage(bot, message.Chat.ID, l["welcome"]) // Отправить приветственное сообщение
	case "poll":
		utils.SendPoll(bot, message.Chat.ID) // Создать опрос
	case "dellnote":
		utils.DeleteNote(bot, message, l) // Создать опрос
	default:
		utils.SendMessage(bot, message.Chat.ID, l["unknown_command"]) // Сообщение о неизвестной команде
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
