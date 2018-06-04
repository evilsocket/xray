package main

import (
	// "fmt"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	/*
		for i := 0; i < 100; i++ {
			func(n int) {
				runner.Run(func() error {
					fmt.Printf("running #%d\n", n)
					return nil
				})
			}(i)
		}
	*/

	in := "google.com"
	intype := units.InputTypeDomain

	for _, u := range units.Loaded {
		if u.AcceptsInputType(intype) {
			u.Run(in)
		}
	}

	runner.Wait()
}
