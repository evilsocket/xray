package units

import "sync"

type State struct {
	sync.RWMutex
	processed map[string]bool
}

func NewState() *State {
	return &State{
		processed: make(map[string]bool),
	}
}

func (s *State) AddProcessed(input string) {
	s.Lock()
	defer s.Unlock()
	s.processed[input] = true
}

func (s *State) DidProcessInput(input string) (found bool) {
	s.RLock()
	defer s.RUnlock()
	_, found = s.processed[input]
	return
}
