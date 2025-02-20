package main

import (
	//"fmt"
	"github.com/joho/godotenv"
	"github.com/net22sky/telegram-bot/bot"
	"github.com/net22sky/telegram-bot/config"
	"github.com/net22sky/telegram-bot/db" // Импортируем пакет db

	"log"
	"os"
	"fmt"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Println("config : ", cfg)
	// Загружаем строки локализации
	locales, err := config.LoadLocales("config/locales.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки строк локализации: %v", err)
	}

	// Установите язык по умолчанию (например, ru)
	lang := "ru"

	// Получаем DSN MySQL из переменной окружения
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		log.Fatal("MYSQL_DSN не установлен")
	}

	// Инициализация GORM
	/*dbInstance, err = db.InitDB(mysqlDSN , true)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}*/

	// Инициализация базы данных
	dbInstance, err := db.InitDB(mysqlDSN, true) // true для автоматических миграций
    if err != nil {
        log.Fatalf("Ошибка при инициализации базы данных: %v", err)
    }

	debug := true
	// Создаем и настраиваем бота
	botInstance, err := bot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"),dbInstance, debug)
	if err != nil {
		log.Panic(err)
	}

	// Запускаем бота
	botInstance.StartPolling( locales, lang)
}
