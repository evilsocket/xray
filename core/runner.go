package core

import (
	"log"
	"runtime"
	"sync"
)

type Job func() error

type Runner struct {
	consumers int
	wg        sync.WaitGroup
	stopped   chan bool
	jobs      chan Job
}

func NewRunner(consumers int) *Runner {
	if consumers <= 0 {
		consumers = runtime.NumCPU() * 2
	}

	return &Runner{
		consumers: consumers,
		wg:        sync.WaitGroup{},
		jobs:      make(chan Job),
		stopped:   make(chan bool),
	}
}

func (r *Runner) worker(id int) {
	// log.Printf("started consumer %d", id)
	for job := range r.jobs {
		if job == nil {
			log.Printf("stopping consumer %d", id)
			r.stopped <- true
			return
		}

		if err := job(); err != nil {
			log.Printf("error while executing Job: %v", err)
		}

		r.wg.Done()
	}
}

func (r *Runner) Start() {
	log.Printf("starting runner with %d consumers", r.consumers)

	for i := 0; i < r.consumers; i++ {
		go r.worker(i)
	}
}

func (r *Runner) Run(j Job) {
	r.wg.Add(1)
	r.jobs <- j
}

func (r *Runner) Wait() {
	r.wg.Wait()
}

func (r *Runner) Stop() {
	for i := 0; i < r.consumers; i++ {
		r.jobs <- nil
	}
	<-r.stopped
}
