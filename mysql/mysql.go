package mysql

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// User представляет информацию о пользователе
type User struct {
    ID         int       `db:"id"`
    TelegramID int64     `db:"telegram_id"`
    Username   string    `db:"username"`
    FirstName  string    `db:"first_name"` // Новое поле для имени пользователя
    Language   string    `db:"language"`
    CreatedAt  time.Time `db:"created_at"`
}

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
    db, err = sql.Open("mysql", dataSourceName+"?parseTime=true")
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }

    log.Println("Подключено к MySQL")

    // Создаем таблицу users
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            telegram_id BIGINT UNIQUE NOT NULL,
            username VARCHAR(255),
            first_name VARCHAR(255),
            language VARCHAR(10) DEFAULT 'ru',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
    if err != nil {
        return err
    }

    // Создаем таблицу notes
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS notes (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id BIGINT NOT NULL,
            text TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(telegram_id) ON DELETE CASCADE
        );
    `)
    if err != nil {
        return err
    }

    return nil
}

// CreateUser создает новую запись о пользователе
func CreateUser(telegramID int64, username string, firstName string) error {
    _, err := db.Exec("INSERT INTO users (telegram_id, username, first_name) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = ?, first_name = ?", 
        telegramID, username, firstName, username, firstName)
    return err
}

// GetUser получает информацию о пользователе по его telegram_id
func GetUser(telegramID int64) (*User, error) {
    row := db.QueryRow("SELECT id, telegram_id, username, first_name, language, created_at FROM users WHERE telegram_id = ?", telegramID)

    var user User
    err := row.Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.Language, &user.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, nil // Пользователь не найден
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// CreateNote создает новую заметку для пользователя
func CreateNote(userID int64, text string) error {
    _, err := db.Exec("INSERT INTO notes (user_id, text) VALUES (?, ?)", userID, text)
    if err != nil {
        log.Printf("Ошибка при создании заметки: %v", err)
        return err
    }
    return nil
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