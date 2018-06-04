package units

import "log"

var (
	Loaded []Unit
)

func init() {
	log.Printf("autoloading units")

	Loaded = make([]Unit, 0)
	Loaded = append(Loaded, NewDNSEnum())
}
