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

// StateManagerInterface определяет методы для управления состояниями.
type StateManagerInterface interface {
	SetUserState(userID int64, newState UserState)
	GetUserState(userID int64) (UserState, bool)
	DeleteUserState(userID int64)
}

// StateManager управляет состояниями пользователей.
type StateManager struct {
	states map[int64]UserState // Хранилище состояний пользователей
	mu     sync.Mutex          // Мьютекс для потокобезопасного доступа
}

// NewStateManager создает новый экземпляр StateManager.
func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[int64]UserState),
	}
}

// SetUserState устанавливает новое состояние для пользователя.
// Параметры:
//   - userID: ID пользователя (Telegram ID)
//   - newState: Новое состояние (см. константы UserState)
func (sm *StateManager) SetUserState(userID int64, newState UserState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.states[userID] = newState
	log.Printf("[State] User %d: state set to '%s'", userID, newState)
}

// GetUserState возвращает текущее состояние пользователя.
// Параметры:
//   - userID: ID пользователя (Telegram ID)
//
// Возвращает:
//   - UserState: Текущее состояние
//   - bool: Флаг наличия состояния (true если состояние существует)
func (sm *StateManager) GetUserState(userID int64) (UserState, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	state, exists := sm.states[userID]
	return state, exists
}

// DeleteUserState удаляет состояние пользователя.
// Параметры:
//   - userID: ID пользователя (Telegram ID)
func (sm *StateManager) DeleteUserState(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, userID)
	log.Printf("[State] User %d: state cleared", userID)
}
