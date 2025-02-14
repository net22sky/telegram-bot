// mysql/mysql.go
package mysql

import (
	"database/sql"
	"encoding/json"
	"github.com/net22sky/telegram-bot/db" // Используем общий пакет db
	"time"
)

// Note представляет заметку пользователя
type Note struct {
	ID        int       `db:"id"`
	UserID    int64     `db:"user_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

// InitMySQL инициализирует подключение к MySQL
func InitMySQL(dataSourceName string) error {
	return db.InitDB(dataSourceName) // Используем db.InitDB
}

// CreateNote создает новую заметку для пользователя
func CreateNote(userID int64, text string) error {
	return db.Exec("INSERT INTO notes (user_id, text) VALUES (?, ?)", userID, text)
}

// GetNotes получает все заметки пользователя
func GetNotes(userID int64) ([]Note, error) {
	rows, err := db.Query("SELECT id, user_id, text, created_at FROM notes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.UserID, &note.Text, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// DeleteNoteByID удаляет заметку по её ID
func DeleteNoteByID(noteID int) error {
	return db.Exec("DELETE FROM notes WHERE id = ?", noteID)
}

// SavePollAnswer сохраняет ответ на опрос в базу данных
func SavePollAnswer(userID int64, pollID string, optionIDs []int) error {
	jsonOptionIDs, err := json.Marshal(optionIDs)
	if err != nil {
		return err
	}

	return db.Exec("INSERT INTO poll_answers (user_id, poll_id, option_ids) VALUES (?, ?, ?)",
		userID, pollID, jsonOptionIDs) // Используем db.Exec
}

// GetNoteByID получает заметку по её ID.
// Параметры:
//   - noteID: ID заметки.
//
// Возвращает:
//   - *Note: Указатель на структуру Note, если заметка найдена.
//   - error: Ошибку, если запрос не удалось выполнить или заметка не существует.
func GetNoteByID(noteID int64) (*Note, error) {
	// Выполняем SQL-запрос
	row, err := db.QueryRow("SELECT id, user_id, text, created_at FROM notes WHERE id = ?", noteID)

	// Создаем временную переменную для хранения данных заметки
	var note Note

	// Сканируем результат запроса
	err = row.Scan(&note.ID, &note.UserID, &note.Text, &note.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Заметка не найдена
	}
	if err != nil {
		return nil, err // Произошла ошибка при сканировании
	}

	return &note, nil // Возвращаем найденную заметку
}
