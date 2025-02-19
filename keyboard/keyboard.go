// keyboard/keyboard.go

package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db"
	"log"
)

// Note представляет заметку пользователя для использования в клавиатуре.
// Это структура используется только для передачи данных между функциями.
type Note struct {
	ID   int    `db:"id"`   // ID заметки
	Text string `db:"text"` // Текст заметки
}

// StartKeyboard создает Inline Keyboard для приветственного сообщения.
// Параметры:
// - Нет параметров.
// Возвращает:
// - tgbotapi.InlineKeyboardMarkup: Клавиатура с кнопками для основных действий.
func StartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заметки", "notes_menu"),    // Кнопка для добавления заметки
			tgbotapi.NewInlineKeyboardButtonData("Напоминания", "reminders_menu"), // Кнопка для удаления заметки
		),
		tgbotapi.NewInlineKeyboardRow(
		 // Кнопка для просмотра списка заметок
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),               // Кнопка для получения справки
		),
	)
}

func NotesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить заметку", "add_note"),    // Кнопка для добавления заметки
			tgbotapi.NewInlineKeyboardButtonData("Удалить заметку", "deletes_note"), // Кнопка для удаления заметки
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список заметок", "view_notes"), // Кнопка для просмотра списка заметок
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),               // Кнопка для получения справки
		),
	)
}

func RemindersKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить напоминание", "add_reminders"),    // Кнопка для добавления заметки
			tgbotapi.NewInlineKeyboardButtonData("Удалить напоминание", "deletes_reminders"), // Кнопка для удаления заметки
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список напоминаний", "view_reminders"), // Кнопка для просмотра списка заметок
			tgbotapi.NewInlineKeyboardButtonData("Помощь", "help"),               // Кнопка для получения справки
		),
	)
}


// DeleteNotesKeyboard создает Inline Keyboard для удаления заметок.
// Параметры:
// - notes []db.Note: Срез заметок пользователя.
// Возвращает:
// - tgbotapi.InlineKeyboardMarkup: Клавиатура с кнопками для каждой заметки.
func DeleteNotesKeyboard(notes []db.Note) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Создаем кнопки для каждой заметки
	for _, note := range notes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d. %s", note.ID, note.Text), fmt.Sprintf("delete_%d", note.ID)),
		))
	}

	// Добавляем кнопку "Отмена"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// LanguageKeyboard создает Inline Keyboard для выбора языка.
// Параметры:
// - Нет параметров.
// Возвращает:
// - tgbotapi.InlineKeyboardMarkup: Клавиатура с кнопками выбора языка.
func LanguageKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "set_language_ru"), // Кнопка для выбора русского языка
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "set_language_en"), // Кнопка для выбора английского языка
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"), // Кнопка для отмены действия
		),
	)
}

// SettingsKeyboard создает Inline Keyboard для настроек пользователя.
// Параметры:
// - Нет параметров.
// Возвращает:
// - tgbotapi.InlineKeyboardMarkup: Клавиатура с кнопками для изменения настроек.
func SettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выбрать язык", "choose_language"),           // Кнопка для выбора языка
			tgbotapi.NewInlineKeyboardButtonData("Настроить часовой пояс", "choose_timezone"), // Кнопка для выбора часового пояса
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"), // Кнопка для отмены действия
		),
	)
}

// CreateNumberKeyboard создает Reply Keyboard с числовыми кнопками.
// Параметры:
// - start int: Начальное число.
// - end int: Конечное число.
// - rowSize int: Количество кнопок в одной строке.
// Возвращает:
// - tgbotapi.ReplyKeyboardMarkup: Клавиатура с числовыми кнопками.
func CreateNumberKeyboard(start, end, rowSize int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	// Создаем ряды кнопок
	for i := start; i <= end; i += rowSize {
		var row []tgbotapi.KeyboardButton

		// Добавляем кнопки в текущую строку
		for j := i; j < i+rowSize && j <= end; j++ {
			row = append(row, tgbotapi.NewKeyboardButton(fmt.Sprintf("%d", j))) // Числовая кнопка
		}

		keyboard = append(keyboard, row) // Добавляем строку в клавиатуру
	}

	return tgbotapi.NewReplyKeyboard(keyboard...) // Создаем Reply Keyboard
}

// CreateMonthKeyboard создает Reply Keyboard с кнопками месяцев.
// Параметры:
// - months []string: Срез названий месяцев.
// Возвращает:
// - tgbotapi.ReplyKeyboardMarkup: Клавиатура с кнопками месяцев.
func CreateMonthKeyboard(months []string) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	// Создаем ряды кнопок по 4 месяца в каждом
	for i := 0; i < len(months); i += 4 {
		end := i + 4
		if end > len(months) {
			end = len(months)
		}

		// Создаем строку с кнопками
		row := make([]tgbotapi.KeyboardButton, 0)
		for _, month := range months[i:end] {
			row = append(row, tgbotapi.NewKeyboardButton(month)) // Кнопка с названием месяца
		}

		keyboard = append(keyboard, row) // Добавляем строку в клавиатуру
	}

	return tgbotapi.NewReplyKeyboard(keyboard...) // Создаем Reply Keyboard
}

// RemoveReplyKeyboard удаляет текущую Reply Keyboard.
// Параметры:
//   - bot: Экземпляр Telegram-бота.
//   - chatID: ID чата.
//   - text: Текст сообщения.
func RemoveReplyKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string) {
	// Создаем конфигурацию для удаления клавиатуры
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // Удаляем клавиатуру

	// Отправляем сообщение
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при удалении клавиатуры: %v", err)
	}
}

func ConfirmationKeyboard() tgbotapi.ReplyKeyboardMarkup {
    keyboard := tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("Сохранить"),
            tgbotapi.NewKeyboardButton("Отменить"),
        ),
    )
// Устанавливаем параметр ResizeKeyboard
keyboard.ResizeKeyboard = true

	return keyboard
}

func ReminderTypesKeyboard() tgbotapi.ReplyKeyboardMarkup {
    keyboard := tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("День рождения"),
            tgbotapi.NewKeyboardButton("Credit"),
            tgbotapi.NewKeyboardButton("Подписка"),
        ),
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("Отменить"),
        ),
    )

    // Устанавливаем параметр ResizeKeyboard
    keyboard.ResizeKeyboard = true

    return keyboard
}