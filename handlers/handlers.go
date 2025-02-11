package handlers

import (
    "log"
    "strings"

    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/yourusername/telegram-bot/mysql"
)

// Locales содержит строки для разных языков
type Locales map[string]map[string]string

// HandleMessage обрабатывает входящие текстовые сообщения
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, locales Locales, lang string) {
    log.Printf("[%s] %s", message.From.UserName, message.Text)

    l := locales[lang] // Выбираем строки для текущего языка

    if strings.HasPrefix(message.Text, "/note ") {
        CreateNote(bot, message, l)
        return
    }

    switch message.Command() {
    case "notes":
        ViewNotes(bot, message, l)
    case "start":
        SendMessage(bot, message.Chat.ID, l["welcome"])
    case "poll":
        SendPoll(bot, message.Chat.ID)
    default:
        SendMessage(bot, message.Chat.ID, l["unknown_command"])
    }
}

// CreateNote создает новую заметку для пользователя
func CreateNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string) {
    parts := strings.SplitN(message.Text, " ", 2)
    if len(parts) < 2 || parts[1] == "" {
        SendMessage(bot, message.Chat.ID, l["unknown_command"])
        return
    }

    noteText := parts[1]
    userID := message.From.ID

    err := mysql.CreateNote(userID, noteText)
    if err != nil {
        log.Printf("Ошибка при создании заметки: %v", err)
        SendMessage(bot, message.Chat.ID, l["note_creation_error"])
        return
    }

    SendMessage(bot, message.Chat.ID, fmt.Sprintf(l["note_created"], noteText))
}

// ViewNotes показывает все заметки пользователя
func ViewNotes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]string) {
    userID := message.From.ID
    notes, err := mysql.GetNotes(userID)
    if err != nil {
        log.Printf("Ошибка при получении заметок: %v", err)
        SendMessage(bot, message.Chat.ID, l["note_retrieval_error"])
        return
    }

    if len(notes) == 0 {
        SendMessage(bot, message.Chat.ID, l["no_notes"])
        return
    }

    var response string
    for i, note := range notes {
        response += fmt.Sprintf("%d. %s (ID: %d)\n", i+1, note.Text, note.ID)
    }

    SendMessage(bot, message.Chat.ID, l["notes_list"]+response)
}