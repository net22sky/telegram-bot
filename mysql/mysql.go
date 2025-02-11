package mysql

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// Note представляет заметку пользователя
type Note struct {
    ID        int       `db:"id"`
    UserID    int64     `db:"user_id"`
    Text      string    `db:"text"`
    CreatedAt time.Time `db:"created_at"`
}

var db *sql.DB

// InitMySQL инициализирует подключение к MySQL
func InitMySQL(dataSourceName string) error {
    var err error
    db, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }

    log.Println("Подключено к MySQL")

    // Создаем таблицу notes, если её нет
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS notes (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id BIGINT NOT NULL,
            text TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
    if err != nil {
        return err
    }

    return nil
}

// CreateNote создает новую заметку для пользователя
func CreateNote(userID int64, text string) error {
    _, err := db.Exec("INSERT INTO notes (user_id, text) VALUES (?, ?)", userID, text)
    return err
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
    _, err := db.Exec("DELETE FROM notes WHERE id = ?", noteID)
    return err
}