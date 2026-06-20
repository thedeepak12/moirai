package pool

import (
	"context"
	"sync"
)

type Pool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
}

func NewPool(numWorkers int, queueSize int) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, queueSize),
		results:    make(chan Result, queueSize),
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 1; i <= p.numWorkers; i++ {
		p.wg.Add(1)

		go func(workerID int) {
			defer p.wg.Done()

			worker(ctx, workerID, p.jobs, p.results)
		}(i)
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Results() <-chan Result {
	return p.results
}

func (p *Pool) Wait() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
