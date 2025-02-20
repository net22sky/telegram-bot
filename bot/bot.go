package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/net22sky/telegram-bot/db/repositories"
	"github.com/net22sky/telegram-bot/db/services"
	"github.com/net22sky/telegram-bot/handlers"
	"github.com/net22sky/telegram-bot/utils"
	"gorm.io/gorm"
)

// Bot содержит зависимости бота.
type Bot struct {
	BotAPI        *tgbotapi.BotAPI
	NoteService   *services.NoteService
	UserService   *services.UserService
	AnswerService *services.PollAnswerService
	Debug         bool
}

// NewBot создает и возвращает новый экземпляр бота
func NewBot(token string, dbInstance *gorm.DB, debug bool) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	botAPI.Debug = debug // Устанавливаем режим отладки

	log.Printf("Авторизован как %s", botAPI.Self.UserName)

	// Создание репозиториев
	noteRepo := repositories.NewNoteRepository(dbInstance)
	userRepo := repositories.NewUserRepository(dbInstance)
	answerRepo := repositories.NewPollAnswerRepository(dbInstance)

	// Создание сервисов
	noteService := services.NewNoteService(noteRepo, userRepo)
	userService := services.NewUserService(userRepo)
	answerService := services.NewPollAnswerService(answerRepo)

	return &Bot{
		BotAPI:        botAPI,
		NoteService:   noteService,
		UserService:   userService,
		AnswerService: answerService,
		Debug:         debug,
	}, nil
}

// SetupMenu настраивает меню команд для бота.
func (b *Bot) SetupMenu() {
	// Создаем список команд
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Получить помощь"},
		{Command: "notes", Description: "Посмотреть список заметок"},

		{Command: "settings", Description: "Настройки языка и времени"},
	}

	// Создаем запрос на установку команд
	setMyCommandsConfig := tgbotapi.NewSetMyCommands(commands...)

	// Отправляем запрос
	_, err := b.BotAPI.Request(setMyCommandsConfig)
	if err != nil {
		log.Panicf("Ошибка при установке меню команд: %v", err)
	}

	log.Println("Меню команд успешно установлено")
}

// StartPolling запускает бота и начинает обработку входящих сообщений.
// Параметры:
//   - bot: Экземпляр Telegram-бота.
//   - locales: Строки локализации.
//   - lang: Язык пользователя.
func (b *Bot) StartPolling(locales map[string]map[string]interface{}, lang string) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.BotAPI.GetUpdatesChan(u)

	b.SetupMenu()

	for update := range updates {
		if update.CallbackQuery != nil {
			handlers.HandleCallbackQuery(b.BotAPI, update.CallbackQuery, locales, lang, b.NoteService, b.UserService)

		}
		if update.Message != nil {
			//handlers.HandleMessage(bot, update.Message, locales, lang)
			handlers.HandleMessage(b.BotAPI, update, locales, lang, b.NoteService, b.UserService)
		}
		if update.PollAnswer != nil {
			utils.HandlePollAnswer(b.BotAPI, update.PollAnswer, b.AnswerService)
		}
	}
}

// SendMessage отправляет сообщение
func (b *Bot) SendMessage(chatID int64, text string, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = replyMarkup
	if _, err := b.BotAPI.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// SendPoll отправляет опрос
func (b *Bot) SendPoll(chatID int64) {
	poll := tgbotapi.NewPoll(chatID, "Какой ваш любимый язык программирования?", "Go", "Python", "JavaScript", "Java")
	poll.IsAnonymous = false // Опрос не анонимный
	if _, err := b.BotAPI.Send(poll); err != nil {
		log.Printf("Ошибка при отправке опроса: %v", err)
	}
}
