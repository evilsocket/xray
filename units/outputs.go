package units

type OutputType int

const (
	OutputTypeDomain OutputType = iota
	OutputTypeIPv4
	OutputTypeIPv6
	OutputTypeIP
	OutputTypeMAC
	OutputTypeOpenPort
)
