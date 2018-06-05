package main

import (
	// "fmt"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"github.com/evilsocket/xray/core"
	"github.com/evilsocket/xray/storage"
	"github.com/evilsocket/xray/units"
)

var (
	input       = flag.String("in", "", "Input value (examples: 'domain:google.com', 'ip:192.168.1.1', ...)")
	consumers   = flag.Int("consumers", 0, "Number of parallel workers or 0 to auto scale.")
	storagePath = flag.String("storage", storage.DefaultPath, "Path of the keys.yml and wordlist files.")
	logFile     = flag.String("log-file", "", "If filled, xray will log to this file.")

	err     = (error)(nil)
	state   = units.NewState()
	runner  = (*core.Runner)(nil)
	strg    = (*storage.Storage)(nil)
	sigChan = (chan os.Signal)(nil)
)

func setupSignals() {
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sigChan
		log.Printf("got signal %v", sig)
		os.Exit(0)
	}()
}

func propagate(input units.Data) []<-chan units.Data {
	chans := make([]<-chan units.Data, 0)

	if !state.DidProcess(input.Data) {
		state.Add(input.Data)
		log.Printf("propagating %s", input)

		for _, u := range units.Loaded {
			if u.AcceptsInput(input) {
				chans = append(chans, u.Run(input))
			}
		}
	}

	return chans
}

func main() {
	flag.Parse()

	setupSignals()

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if strg, err = storage.Open(*storagePath, true); err != nil {
		log.Fatalf("error opening storage: %v", err)
	}

	runner = core.NewRunner(*consumers)
	runner.Start()
	defer runner.Stop()

	in := units.Data{
		Type: units.DataTypeDomain,
		Data: "something-something.ansa.it",
	}

	queue := []units.Data{in}

	for {
		// log.Printf("%v", queue)
		for _, input := range queue {
			chans := propagate(input)

			for _, ch := range chans {
				func(c <-chan units.Data) {
					go func() {
						for out := range c {
							log.Printf("  > %s", out)
							for _, o := range out.Explode() {
								queue = append(queue, o)
							}
						}
					}()
				}(ch)
			}
		}

		// time.Sleep(100 * time.Millisecond)
		runner.Wait()
	}
}
