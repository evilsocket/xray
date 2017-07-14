/*
 * Copyleft 2017, Simone Margaritelli <evilsocket at protonmail dot com>
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *   * Neither the name of ARM Inject nor the names of its contributors may be used
 *     to endorse or promote products derived from this software without
 *     specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */
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
