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
