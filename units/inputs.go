package units

type InputType int

const (
	InputTypeDomain InputType = iota
	InputTypeIPv4
	InputTypeIPv6
	InputTypeIP
	InputTypeMAC
	InputTypeOpenPort
)
