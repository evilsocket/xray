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
	"github.com/ns3777k/go-shodan/shodan"
	"sort"
	"sync"
)

type HistoryEntry struct {
	Address  string `json:"ip"`
	Location string `json:"location"`
	ISP      string `json:"owner"`
	Updated  string `json:"lastseen"`
}

type Target struct {
	Address string
	Domains []string
	Banners map[string]string
	Info    *shodan.Host
	History map[string][]HistoryEntry

	ctx      *Context
	grabbers []Grabber
	lock     sync.Mutex
}

func NewTarget(address string, domain string) *Target {
	t := &Target{
		Address:  address,
		Domains:  []string{domain},
		Banners:  make(map[string]string),
		History:  make(map[string][]HistoryEntry),
		Info:     nil,
		grabbers: make([]Grabber, 0),
		ctx:      GetContext(),
	}

	t.grabbers = append(t.grabbers, &HTTPGrabber{})
	t.grabbers = append(t.grabbers, &DNSGrabber{})
	t.grabbers = append(t.grabbers, &MYSQLGrabber{})
	t.grabbers = append(t.grabbers, NewLineGrabber("smtp", []int{25, 587}))
	t.grabbers = append(t.grabbers, NewLineGrabber("ftp", []int{21}))
	t.grabbers = append(t.grabbers, NewLineGrabber("ssh", []int{22, 222, 2222}))
	t.grabbers = append(t.grabbers, NewLineGrabber("pop", []int{110}))
	t.grabbers = append(t.grabbers, NewLineGrabber("irc", []int{6667}))

	t.scanDomainAsync(domain)
	t.startAsyncScan()
	return t
}

func (t Target) AddDomain(domain string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	for _, d := range t.Domains {
		if d == domain {
			return false
		}
	}

	t.Domains = append(t.Domains, domain)

	if _, ok := t.History[domain]; ok == false {
		t.scanDomainAsync(domain)
	}

	return true
}

func (t *Target) SortedBanners() []string {
	banners := make([]string, 0, len(t.Banners))
	for name := range t.Banners {
		banners = append(banners, name)
	}
	sort.Strings(banners)
	return banners
}

func (t Target) scanDomainAsync(domain string) {
	go func(t *Target, domain string) {
		t.lock.Lock()
		defer t.lock.Unlock()
		t.History[domain] = t.ctx.VDNS.GetHistory(domain)
	}(&t, domain)
}

func (t *Target) startAsyncScan() {
	go func() {
		info, err := t.ctx.Shodan.GetServicesForHost(t.Address, &shodan.HostServicesOptions{
			History: false,
			Minify:  true,
		})
		if err == nil {
			t.Info = info
			go t.startAsyncBannerGrabbing()
		}
	}()
}

func (t *Target) startAsyncBannerGrabbing() {
	go func() {
		if t.Info != nil {
			for _, port := range t.Info.Ports {
				for _, grabber := range t.grabbers {
					grabber.Grab(port, t)
				}
			}
		}
	}()
}
