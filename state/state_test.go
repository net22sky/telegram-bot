package state

import (
	"sync"
	"testing"
)

func TestConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			SetUserState(id, StateAddingNote)
			GetUserState(id)
			DeleteUserState(id)
		}(int64(i))
	}
	wg.Wait()
}
