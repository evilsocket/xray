package units

type DataChannel <-chan Data

type Unit interface {
	AcceptsInput(in Data) bool
	Propagates() bool
	Run(in Data) <-chan Data
}
