package state

import (
	"log"
	"sync"
)

var userStates = make(map[int64]string) // Хранилище состояний пользователей
var userStatesMutex = &sync.Mutex{}     // Мьютекс для защиты userStates

const (
	StateIdle       = "idle"        // Режим ожидания
	StateAddingNote = "adding_note" // Режим добавления заметки
)

// SetUserState устанавливает новое состояние пользователя.
// Параметры:
//   - userID: ID пользователя.
//   - newState: Новое состояние (например, "adding_note").
func SetUserState(userID int64, newState string) {
	userStatesMutex.Lock()
	defer userStatesMutex.Unlock()

	userStates[userID] = newState
	log.Printf("Пользователь %d: Состояние установлено как %s", userID, newState)
}

// GetUserState получает текущее состояние пользователя.
// Параметры:
//   - userID: ID пользователя.
//
// Возвращает:
//   - string: Текущее состояние пользователя.
//   - bool: true, если состояние существует; false, если нет.
func GetUserState(userID int64) (string, bool) {
	userStatesMutex.Lock()
	defer userStatesMutex.Unlock()

	state, exists := userStates[userID]
	return state, exists
}

// DeleteUserState удаляет состояние пользователя.
// Параметры:
//   - userID: ID пользователя.
func DeleteUserState(userID int64) {
	userStatesMutex.Lock()
	defer userStatesMutex.Unlock()

	delete(userStates, userID)
	log.Printf("Пользователь %d: Состояние очищено", userID)
}
