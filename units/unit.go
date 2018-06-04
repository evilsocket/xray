package units

type Unit interface {
	AcceptsInputType(t InputType) bool
	EmitsOutputType() OutputType
	Propagates() bool

	Run(input string)
}
