package core

import (
	"log"
	"runtime"
	"sync"
)

type Job func() error

type Runner struct {
	name      string
	consumers int
	wg        sync.WaitGroup
	stopped   chan bool
	jobs      chan Job
}

func NewRunner(consumers int, name string) *Runner {
	if consumers <= 0 {
		consumers = runtime.NumCPU() * 2
	}

	return &Runner{
		consumers: consumers,
		name:      name,
		wg:        sync.WaitGroup{},
		jobs:      make(chan Job),
		stopped:   make(chan bool),
	}
}

func (r *Runner) worker(id int) {
	// log.Printf("started consumer %d", id)
	for job := range r.jobs {
		if job == nil {
			log.Printf("%s: stopping consumer %d", r.name, id)
			r.stopped <- true
			return
		}

		if err := job(); err != nil {
			log.Printf("%s: error while executing Job: %v", r.name, err)
		}

		r.wg.Done()
	}
}

func (r *Runner) Start() {
	log.Printf("%s: starting runner with %d consumers", r.name, r.consumers)
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
