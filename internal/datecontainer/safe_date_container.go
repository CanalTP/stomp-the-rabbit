package datecontainer

import (
	"sync"
	"time"
)

type SafeDateContainer struct {
	mu     sync.Mutex
	values map[string]time.Time
}

func NewSafeDateContainer() (sb SafeDateContainer) {
	sb.values = make(map[string]time.Time)
	return
}

// Safely get the time of the given counter
func (s *SafeDateContainer) Get(key string) time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values[key]
}

// Safely set the given counter to now
func (s *SafeDateContainer) Refresh(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = time.Now()
}

// Safely get all the times
func (s *SafeDateContainer) GetAll() map[string]time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values
}
