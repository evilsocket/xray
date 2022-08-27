package xray

import (
	"sort"
	"sync"
)

type Pool struct {
	sync.RWMutex
	addrs   []string
	Session *Session
}

func NewPool(session *Session) *Pool {
	return &Pool{
		Session: session,
		addrs:   make([]string, 0),
	}
}

func (p *Pool) WasRestored() bool {
	return len(p.Session.Targets) > 0
}

func (p *Pool) FlushSession(stats *Statistics) {
	p.Lock()
	defer p.Unlock()

	p.Session.Flush(stats)
}

func (p *Pool) Find(address string) *Target {
	p.RLock()
	defer p.RUnlock()

	t, found := p.Session.Targets[address]
	if found {
		return t
	} else {
		return nil
	}
}

func (p *Pool) Add(t *Target) {
	p.Lock()
	defer p.Unlock()

	p.Session.Targets[t.Address] = t
}

func (p *Pool) Sorted() []string {
	if len(p.addrs) == 0 {
		p.addrs = make([]string, 0, len(p.Session.Targets))
		for addr := range p.Session.Targets {
			p.addrs = append(p.addrs, addr)
		}
		sort.Strings(p.addrs)
	}

	return p.addrs
}
