package units

import (
	"fmt"
	"log"
	"net"

	"github.com/evilsocket/xray/core"
	"github.com/evilsocket/xray/storage"

	"github.com/bobesa/go-domain-util/domainutil"
)

type DNSEnum struct {
	state  *State
	output chan Data
}

func NewDNSEnum() *DNSEnum {
	return &DNSEnum{
		state:  NewState(),
		output: make(chan Data),
	}
}

func (d DNSEnum) AcceptsInput(in Data) bool {
	return in.Type == DataTypeDomain
}

func (d DNSEnum) Propagates() bool {
	return true
}

func (d DNSEnum) Run(in Data) <-chan Data {
	domain := domainutil.Domain(in.Data)
	if d.state.DidProcess(domain) == false {
		d.state.Add(domain)

		log.Printf("dns:enum(%s)", domain)

		go func() {
			for _, word := range storage.I.Domains {
				// save context
				func(subdomain string) {
					core.Queue.Run(func() error {
						hostname := fmt.Sprintf("%s.%s", subdomain, domain)
						if addrs, err := net.LookupHost(hostname); err == nil {
							d.output <- Data{
								Type:  DataTypeDomain,
								Data:  hostname,
								Extra: addrs,
							}
						}
						return nil
					})
				}(word)
			}
		}()
	}

	return d.output
}
