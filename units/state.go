package units

import (
	"sync"
)

type State struct {
	sync.RWMutex
	processed map[string]bool
}

func NewState() *State {
	return &State{
		processed: make(map[string]bool),
	}
}

func (s *State) Add(data string) {
	s.Lock()
	defer s.Unlock()
	s.processed[data] = true
}

func (s *State) DidProcess(data string) (found bool) {
	s.RLock()
	defer s.RUnlock()
	_, found = s.processed[data]
	return
}
