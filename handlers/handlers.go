package handlers

import (
    "log"
    "strings"
    "fmt"

    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
   
    "github.com/net22sky/telegram-bot/mysql"
    
)

// Locales содержит строки для разных языков
type Locales map[string]map[string]string

// HandleMessage обрабатывает входящие текстовые сообщения
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, locales Locales, defaultLang string) {
    user := getUser(message.From.ID)
    if user == nil {
        log.Printf("Пользователь %d не найден в базе данных", message.From.ID)
        userLang := defaultLang
        sendDefaultMessage(bot, message.Chat.ID, "Произошла ошибка при получении вашего профиля.", userLang)
        return
    }

    l := locales[user.Language] // Выбираем строки для текущего языка

    log.Printf("[%s] %s", message.From.UserName, message.Text)

    if strings.HasPrefix(message.Text, "/note ") {
        CreateNote(bot, message, l, user.FirstName)
        return
    }

    switch message.Command() {
    case "notes":
        ViewNotes(bot, message, l, user.FirstName)
    case "start":
        sendWelcomeMessage(bot, message.Chat.ID, user.FirstName, l)
    case "poll":
        SendPoll(bot, message.Chat.ID)
    default:
        sendDefaultMessage(bot, message.Chat.ID, l["unknown_command"], user.Language)
    }
}

// getUser получает информацию о пользователе из базы данных
func getUser(telegramID int64) *mysql.User {
    user, err := mysql.GetUser(telegramID)
    if err != nil {
        log.Printf("Ошибка при получении пользователя: %v", err)
        return nil
    }
    return user
}

// sendWelcomeMessage отправляет приветственное сообщение с именем пользователя
func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64, firstName string, l map[string]string) {
    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(l["welcome"], firstName))
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Ошибка при отправке сообщения: %v", err)
    }
}

// sendDefaultMessage отправляет сообщение по умолчанию
func sendDefaultMessage(bot *tgbotapi.BotAPI, chatID int64, text string, lang string) {
    msg := tgbotapi.NewMessage(chatID, text)
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Ошибка при отправке сообщения: %v", err)
    }
}

// CreateNote создает новую заметку для пользователя
func CreateNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string, firstName string) {
    parts := strings.SplitN(message.Text, " ", 2)
    if len(parts) < 2 || parts[1] == "" {
        sendDefaultMessage(bot, message.Chat.ID, l["unknown_command"], l["language"])
        return
    }

    noteText := parts[1]
    userID := message.From.ID

    err := mysql.CreateNote(userID, noteText)
    if err != nil {
        log.Printf("Ошибка при создании заметки: %v", err)
        sendDefaultMessage(bot, message.Chat.ID, l["note_creation_error"], l["language"])
        return
    }

    sendDefaultMessage(bot, message.Chat.ID, fmt.Sprintf(l["note_created"], noteText), l["language"])
}

// ViewNotes показывает все заметки пользователя
func ViewNotes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string, firstName string) {
    userID := message.From.ID
    notes, err := mysql.GetNotes(userID)
    if err != nil {
        log.Printf("Ошибка при получении заметок: %v", err)
        sendDefaultMessage(bot, message.Chat.ID, l["note_retrieval_error"], l["language"])
        return
    }

    if len(notes) == 0 {
        sendDefaultMessage(bot, message.Chat.ID, l["no_notes"], l["language"])
        return
    }

    var response string
    for i, note := range notes {
        response += fmt.Sprintf("%d. %s (ID: %d)\n", i+1, note.Text, note.ID)
    }

    sendDefaultMessage(bot, message.Chat.ID, l["notes_list"]+response, l["language"])
}