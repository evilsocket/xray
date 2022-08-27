package grabbers

import (
	"fmt"
	"regexp"

	"github.com/evilsocket/xray"
	"github.com/miekg/dns"
)

type DNSGrabber struct {
}

func (g *DNSGrabber) Name() string {
	return "dns"
}

func (g *DNSGrabber) grabChaos(addr string, q string) string {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{q, dns.TypeTXT, dns.ClassCHAOS}

	in, _, _ := c.Exchange(m, addr)
	if in != nil && len(in.Answer) > 0 {
		s := in.Answer[0].String()
		re := regexp.MustCompile(".*\"([^\"]+)\".*")
		match := re.FindStringSubmatch(s)
		if len(match) > 0 {
			return match[1]
		}
	}
	return ""
}

func (g *DNSGrabber) Grab(port int, t *xray.Target) {
	if port != 53 {
		return
	}

	addr := fmt.Sprintf("%s:53", t.Address)

	if v := g.grabChaos(addr, "version.bind."); v != "" {
		t.Banners["dns:version"] = v
	}

	if h := g.grabChaos(addr, "hostname.bind."); h != "" {
		t.Banners["dns:hostname"] = h
	}
}
