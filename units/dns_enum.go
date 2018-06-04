package units

import "log"

type DNSEnum struct {
	state *State
}

func NewDNSEnum() *DNSEnum {
	return &DNSEnum{
		state: NewState(),
	}
}

func (d DNSEnum) AcceptsInputType(t InputType) bool {
	return t == InputTypeDomain
}

func (d DNSEnum) EmitsOutputType() OutputType {
	return OutputTypeDomain
}

func (d DNSEnum) Propagates() bool {
	return true
}

func (d DNSEnum) Run(input string) {
	if d.state.DidProcessInput(input) == false {
		d.state.AddProcessed(input)
		log.Printf("dns:enum(%v)", input)
	}
}
