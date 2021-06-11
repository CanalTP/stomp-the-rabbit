package scoreboard

import (
	"sync"
	"time"
)

type ScoreBoard struct {
	mu     sync.Mutex
	values map[string]time.Time
}

func NewScoreBoard() (sb ScoreBoard) {
	sb.values = make(map[string]time.Time)
	return
}

// Safely get the value of the given counter
func (s *ScoreBoard) Get(key string) time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values[key]
}

// Safely increase the value of the given counter
func (s *ScoreBoard) Set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = time.Now()
}

// All the counters
func (s *ScoreBoard) All() map[string]time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values
}
