package state

import (
	"log"
	"sync"
)

// UserState представляет состояние пользователя в приложении.
type UserState string

const (
	StateIdle        UserState = "idle"        // Пользователь в режиме ожидания
	StateAddingNote  UserState = "adding_note" // Пользователь добавляет заметку
	StateEditingNote UserState = "editing_note"
)

var (
	userStates      = make(map[int64]UserState) // Хранилище состояний пользователей
	userStatesMutex = &sync.Mutex{}             // Мьютекс для потокобезопасного доступа
)

// SetUserState устанавливает новое состояние для пользователя.
// Параметры:
//   - userID: ID пользователя (Telegram ID)
//   - newState: Новое состояние (см. константы UserState)
func SetUserState(userID int64, newState UserState) {
	userStatesMutex.Lock()
	defer userStatesMutex.Unlock()

	userStates[userID] = newState
	log.Printf("[State] User %d: state set to '%s'", userID, newState)
}

// GetUserState возвращает текущее состояние пользователя.
// Параметры:
//   - userID: ID пользователя (Telegram ID)
//
// Возвращает:
//   - UserState: Текущее состояние
//   - bool: Флаг наличия состояния (true если состояние существует)
func GetUserState(userID int64) (UserState, bool) {
	userStatesMutex.Lock()
	defer userStatesMutex.Unlock()

	state, exists := userStates[userID]
	return state, exists
}

// DeleteUserState удаляет состояние пользователя.
// Параметры:
//   - userID: ID пользователя (Telegram ID)
func DeleteUserState(userID int64) {
	userStatesMutex.Lock()
	defer userStatesMutex.Unlock()

	delete(userStates, userID)
	log.Printf("[State] User %d: state cleared", userID)
}
