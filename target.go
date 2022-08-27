package xray

import (
	"net"
	"sort"
	"sync"

	"github.com/ns3777k/go-shodan/shodan"
)

type HistoryEntry struct {
	Address  string `json:"ip"`
	Location string `json:"location"`
	ISP      string `json:"owner"`
	Updated  string `json:"lastseen"`
}

type Target struct {
	sync.Mutex

	Address   string
	Hostnames []string
	Domains   []string
	Banners   map[string]string
	Info      *shodan.Host
	History   map[string][]HistoryEntry

	ctx *Context `json:"-"`
}

func NewTarget(address string, domain string) *Target {
	t := &Target{
		Address:   address,
		Hostnames: make([]string, 0),
		Domains:   []string{domain},
		Banners:   make(map[string]string),
		History:   make(map[string][]HistoryEntry),
		Info:      nil,
		ctx:       GetContext(),
	}

	t.scanDomainAsync(domain)
	t.startAsyncScan()
	return t
}

func (t *Target) AddDomain(domain string) bool {
	t.Lock()
	defer t.Unlock()

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

func (t *Target) scanDomainAsync(domain string) {
	go func(t *Target, domain string) {
		t.Lock()
		defer t.Unlock()
		t.History[domain] = t.ctx.VDNS.GetHistory(domain)
	}(t, domain)
}

func (t *Target) startAsyncScan() {
	go func() {
		if names, err := net.LookupAddr(t.Address); err == nil {
			t.Hostnames = names
		}

		info, err := t.ctx.Shodan.GetServicesForHost(nil, t.Address, &shodan.HostServicesOptions{
			History: false,
			Minify:  true,
		})
		if err == nil {
			t.Info = info
			t.ctx.StartGrabbing(t)
		}
	}()
}
