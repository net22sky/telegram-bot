package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db"

	"log"
	"strconv"
	"strings"
)

// SendMessage отправляет текстовое сообщение в указанный чат.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - chatID: ID чата, в который нужно отправить сообщение.
// - text: Текст сообщения.
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text) // Создание нового сообщения
	if _, err := bot.Send(msg); err != nil { // Отправка сообщения

		log.Printf("Ошибка при отправке сообщения: %v", err) // Логирование ошибки, если отправка не удалась
	}
	log.Printf("Ошибка при отправке сообщения: %v", msg) // Логирование ошибки, если отправка не удалась
}

// SendPoll отправляет опрос в указанный чат.
// Параметры:
// - bot: Экземпляр Telegram-бота.
// - chatID: ID чата, в который нужно отправить опрос.
func SendPoll(bot *tgbotapi.BotAPI, chatID int64) {
	// Создание нового опроса с вопросом и вариантами ответов
	poll := tgbotapi.NewPoll(chatID, "Какой ваш любимый язык программирования?", "Go", "Python", "JavaScript", "Java")
	poll.IsAnonymous = false                  // Установка опроса как неанонимного
	if _, err := bot.Send(poll); err != nil { // Отправка опроса
		log.Printf("Ошибка при отправке опроса: %v", err) // Логирование ошибки, если отправка не удалась
	}
}

// DeleteNote удаляет заметку по её ID, если она принадлежит указанному пользователю.
// Параметры:
//   - bot: Экземпляр Telegram-бота.
//   - message: Сообщение от пользователя.
//   - l: Строки локализации для текущего языка.
func DeleteNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}) {
	parts := strings.SplitN(message.Text, " ", 2)
	if len(parts) < 2 || parts[1] == "" {
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "unknown_command"))
		return
	}

	// Извлекаем ID заметки из сообщения
	noteIDStr := parts[1]
	userID := message.From.ID

	// Преобразуем noteID в числовой формат
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil || noteID <= 0 {
		log.Printf("Ошибка при преобразовании ID заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "invalid_note_id"))
		return
	}

	// Получаем заметку из базы данных
	note, err := db.GetNoteByID(int64(noteID))
	if err != nil {
		log.Printf("Ошибка при получении заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_retrieval_error"))
		return
	}

	// Получаем пользователя по Telegram ID
	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Printf("ошибка при получении пользователя: %v", err)
		return
	}

	// Проверяем, что заметка существует и принадлежит текущему пользователю
	if note == nil || note.UserID != uint(user.ID) {
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_not_found"))
		return
	}

	// Удаляем заметку
	err = db.DeleteNoteByID(uint(noteID), int64(userID))
	if err != nil {
		log.Printf("Ошибка при удалении заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_deletion_error"))
		return
	}

	// Отправляем сообщение об успешном удалении
	SendMessage(bot, message.Chat.ID, fmt.Sprintf(GetLocalizedString(l, "note_deleted"), noteID))
}

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}) {
	chatID := message.Chat.ID

	// Отправляем сообщение с помощью локализованного текста
	SendMessage(bot, chatID, GetLocalizedString(l, "help_message"))
}

// getLocalizedString получает локализованную строку по ключу.
func GetLocalizedString(l map[string]interface{}, key string) string {
	if value, exists := l[key]; exists {
		if strValue, ok := value.(string); ok {
			return strValue
		}
		log.Printf("Ошибка преобразования строки для ключа %s", key)
	}
	return "Строка не найдена"
}
