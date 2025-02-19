// db/db_test.go

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB() (*gorm.DB, func()) {
	// Создаем базу данных SQLite в памяти
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Не удалось инициализировать тестовую базу данных")
	}

	// Выполняем миграцию таблиц
	err = db.AutoMigrate(&User{}, &Note{}, &PollAnswer{})
	if err != nil {
		panic("Ошибка при миграции таблиц")
	}

	// Функция для очистки ресурсов после теста
	cleanup := func() {
		db.Exec("DROP TABLE IF EXISTS notes")
		db.Exec("DROP TABLE IF EXISTS users")
		db.Exec("DROP TABLE IF EXISTS poll_answers")
	}

	return db, cleanup
}

func TestCreateUser(t *testing.T) {
	// Инициализируем тестовую базу данных
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Заменяем глобальную переменную DB на тестовую
	DB = testDB

	// Создаем пользователя
	user, err := CreateUser(12345, "testuser", "Test Name")
	assert.NoError(t, err, "Ожидалось успешное создание пользователя")
	assert.NotNil(t, user, "Пользователь должен быть создан")

	// Проверяем, что пользователь существует в базе данных
	foundUser, err := GetUserByID(12345)
	assert.NoError(t, err, "Ожидалось успешное получение пользователя")
	assert.NotNil(t, foundUser, "Пользователь должен существовать")
	assert.Equal(t, uint(1), foundUser.ID, "ID пользователя должно быть 1")
	assert.Equal(t, int64(12345), foundUser.TelegramID, "TelegramID пользователя должно совпадать")
}

func TestCreateAndGetNotes(t *testing.T) {
	// Инициализируем тестовую базу данных
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Заменяем глобальную переменную DB на тестовую
	DB = testDB

	// Создаем пользователя
	user, err := CreateUser(12345, "testuser", "Test Name")
	assert.NoError(t, err, "Ожидалось успешное создание пользователя")
	assert.NotNil(t, user, "Пользователь должен быть создан")

	// Создаем заметку для пользователя
	err = CreateNote(12345, "Купить хлеб")
	assert.NoError(t, err, "Ожидалось успешное создание заметки")

	// Получаем список заметок пользователя
	notes, err := GetNotes(12345)
	assert.NoError(t, err, "Ожидалось успешное получение заметок")
	assert.Len(t, notes, 1, "Должна быть создана одна заметка")

	// Проверяем данные заметки
	note := notes[0]
	assert.Equal(t, uint(1), note.UserID, "UserID заметки должно совпадать с ID пользователя")
	assert.Equal(t, "Купить хлеб", note.Text, "Текст заметки должен совпадать")
}

func TestDeleteNote(t *testing.T) {
	// Инициализируем тестовую базу данных
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Заменяем глобальную переменную DB на тестовую
	DB = testDB

	// Создаем пользователя
	user, err := CreateUser(12345, "testuser", "Test Name")
	assert.NoError(t, err, "Ожидалось успешное создание пользователя")
	assert.NotNil(t, user, "Пользователь должен быть создан")

	// Создаем заметку для пользователя
	err = CreateNote(12345, "Купить хлеб")
	assert.NoError(t, err, "Ожидалось успешное создание заметки")

	// Удаляем заметку
	err = DeleteNoteByID(1, 12345)
	assert.NoError(t, err, "Ожидалось успешное удаление заметки")

	// Проверяем, что заметка удалена
	notes, err := GetNotes(12345)
	assert.NoError(t, err, "Ожидалось успешное получение заметок")
	assert.Len(t, notes, 0, "Заметок должно быть 0 после удаления")
}

func TestSavePollAnswer(t *testing.T) {
	// Инициализируем тестовую базу данных
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Заменяем глобальную переменную DB на тестовую
	DB = testDB

	// Создаем пользователя
	user, err := CreateUser(12345, "testuser", "Test Name")
	assert.NoError(t, err, "Ожидалось успешное создание пользователя")
	assert.NotNil(t, user, "Пользователь должен быть создан")

	// Сохраняем ответ на опрос
	err = SavePollAnswer(uint(user.ID), "poll_123", []int{0, 1})
	assert.NoError(t, err, "Ожидалось успешное сохранение ответа")

	// Проверяем, что ответ существует в базе данных
	var pollAnswers []PollAnswer
	result := DB.Where("user_id = ? AND poll_id = ?", user.ID, "poll_123").Find(&pollAnswers)
	assert.NoError(t, result.Error, "Ожидалось успешное получение ответа")
	assert.Len(t, pollAnswers, 1, "Должен быть сохранен один ответ")
	assert.Equal(t, []int{0, 1}, pollAnswers[0].OptionIDs, "OptionIDs должны совпадать")
}

func TestGetNoteByID(t *testing.T) {
	// Инициализируем тестовую базу данных
	testDB, cleanup := setupTestDB()
	defer cleanup()

	// Заменяем глобальную переменную DB на тестовую
	DB = testDB

	// Создаем пользователя
	user, err := CreateUser(12345, "testuser", "Test Name")
	assert.NoError(t, err, "Ожидалось успешное создание пользователя")
	assert.NotNil(t, user, "Пользователь должен быть создан")

	// Создаем заметку для пользователя
	err = CreateNote(12345, "Купить хлеб")
	assert.NoError(t, err, "Ожидалось успешное создание заметки")

	// Получаем заметку по её ID
	note, err := GetNoteByID(1)
	assert.NoError(t, err, "Ожидалось успешное получение заметки")
	assert.NotNil(t, note, "Заметка должна существовать")
	assert.Equal(t, uint(1), note.ID, "ID заметки должно быть 1")
	assert.Equal(t, "Купить хлеб", note.Text, "Текст заметки должен совпадать")
}
