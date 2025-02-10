package main

import (
    "log"
    "os"

    "github.com/net22sky/telegram-bot/app/config"
    "github.com/net22sky/telegram-bot/app/bot"
    "github.com/net22sky/telegram-bot/app/mongo"
)

func main() {
    // Загружаем конфигурацию
    cfg, err := config.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Ошибка загрузки конфигурации: %v", err)
    }

    // Загружаем строки локализации
    locales, err := config.LoadLocales("config/locales.yaml")
    if err != nil {
        log.Fatalf("Ошибка загрузки строк локализации: %v", err)
    }

    // Установите язык по умолчанию (например, ru)
    lang := "ru"

    // Получаем URI MongoDB из переменной окружения
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI не установлен")
    }

    // Инициализация MongoDB
    err = mongo.InitMongoDB(mongoURI, "telegram_bot", "notes")
    if err != nil {
        log.Fatalf("Ошибка подключения к MongoDB: %v", err)
    }

    // Создаем и настраиваем бота
    tgBot, err := bot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"), cfg.Telegram.Debug)
    if err != nil {
        log.Panic(err)
    }

    // Запускаем бота
    bot.StartPolling(tgBot, locales, lang)
}