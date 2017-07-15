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
