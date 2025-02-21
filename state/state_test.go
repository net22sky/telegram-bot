package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateManager(t *testing.T) {
	manager := NewStateManager()
	userID := int64(123456789)

	// Устанавливаем состояние
	manager.SetUserState(userID, StateAddingNote)

	// Проверяем состояние
	state, exists := manager.GetUserState(userID)
	assert.True(t, exists)
	assert.Equal(t, StateAddingNote, state)

	// Удаляем состояние
	manager.DeleteUserState(userID)

	// Проверяем, что состояние удалено
	_, exists = manager.GetUserState(userID)
	assert.False(t, exists)
}
