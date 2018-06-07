package units

import (
	"log"

	"github.com/evilsocket/xray/core"
)

type PortScan struct {
	state  *State
	runner *core.Runner
	output chan Data
}

func NewPortScan() *PortScan {
	d := &PortScan{
		state:  NewState(),
		runner: core.NewRunner(4, "port:scan::runner"),
		output: make(chan Data),
	}
	d.runner.Start()

	return d
}

func (d PortScan) AcceptsInput(in Data) bool {
	return in.Type == DataTypeIP
}

func (d PortScan) Propagates() bool {
	return true
}

func (d PortScan) Run(in Data) <-chan Data {
	address := in.Data
	if d.state.DidProcess(address) == false {
		d.state.Add(address)

		log.Printf("ip:scan(%s)", address)

		/*
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
		*/
	}

	return d.output
}
