package internal

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workerCount int
	jobs        chan Job
	results     chan Result
	Done        chan struct{}
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, result chan<- Result) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			result <- job.execute(ctx)
		case <-ctx.Done():
			result <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}

func New(wc int) WorkerPool {
	return WorkerPool{
		workerCount: wc,
		jobs:        make(chan Job),
		results:     make(chan Result),
		Done:        make(chan struct{}),
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workerCount; i++ {
		wg.Add(1)

		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp WorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp WorkerPool) GenerateJobBulk(jobBulk []Job) {
	for _, job := range jobBulk {
		wp.jobs <- job
	}

	close(wp.jobs)
}
