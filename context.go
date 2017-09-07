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
	"strings"
	"sync"

	"github.com/ns3777k/go-shodan/shodan"
)

var (
	instance *Context = nil
	lock     sync.Mutex
	Grabbers []Grabber
)

type Grabber interface {
	Name() string
	Grab(port int, t *Target)
}

type Context struct {
	Domain  string
	Bruter  *Machine
	Session *Session
	Pool    *Pool
	Shodan  *shodan.Client
	VDNS    *ViewDNS
	CSH     *CertSH
}

func MakeContext(domain string, session_file string, consumers int, wordlist string, shodan_token string, viewdns_token string, run_handler RunHandler, res_handler ResultHandler) *Context {
	lock.Lock()
	defer lock.Unlock()

	instance = &Context{}
	instance.Domain = domain
	instance.Session = NewSession(session_file)
	instance.Pool = NewPool(instance.Session)
	instance.Bruter = NewMachine(consumers, wordlist, instance.Session, run_handler, res_handler)
	instance.Shodan = shodan.NewClient(nil, shodan_token)
	instance.VDNS = NewViewDNS(viewdns_token)
	instance.CSH = NewCertSH()

	return instance
}

func SetupGrabbers(gs []Grabber) {
	Grabbers = gs
}

func GetContext() *Context {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		// this should not happen as the instance is
		// initialized in main
		panic("(╯°□°）╯︵ ┻━┻")
	}

	return instance
}

func (c *Context) GetSubDomain(domain string) string {
	if strings.HasSuffix(domain, c.Domain) == true && domain != c.Domain {
		subdomain := strings.Replace(domain, "."+c.Domain, "", -1)
		if subdomain != "*" && subdomain != "" {
			return subdomain
		}
	}
	return ""
}

func (c *Context) StartGrabbing(t *Target) {
	go func() {
		if t.Info != nil {
			for _, port := range t.Info.Ports {
				for _, grabber := range Grabbers {
					grabber.Grab(port, t)
				}
			}
		}
	}()
}
