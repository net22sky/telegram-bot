package bot

import (
    "context"
    "log"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/net22sky/telegram-bot/app/handlers"
)

// NewBot создает и возвращает новый экземпляр бота
func NewBot(token string, debug bool) (*tgbotapi.BotAPI, error) {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }
    bot.Debug = debug
    return bot, nil
}

// StartPolling запускает бота и начинает обработку входящих сообщений
func StartPolling(bot *tgbotapi.BotAPI, locales handlers.Locales, lang string) {
    log.Printf("Авторизован как %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            handlers.HandleMessage(bot, update.Message, locales, lang)
        }
        if update.PollAnswer != nil {
            handlers.HandlePollAnswer(bot, update.PollAnswer)
        }
    }
}

// SendMessage отправляет сообщение
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
    msg := tgbotapi.NewMessage(chatID, text)
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Ошибка при отправке сообщения: %v", err)
    }
}

// SendPoll отправляет опрос
func SendPoll(bot *tgbotapi.BotAPI, chatID int64) {
    poll := tgbotapi.NewPoll(chatID, "Какой ваш любимый язык программирования?", "Go", "Python", "JavaScript", "Java")
    poll.IsAnonymous = false // Опрос не анонимный
    if _, err := bot.Send(poll); err != nil {
        log.Printf("Ошибка при отправке опроса: %v", err)
    }
}