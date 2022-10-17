package internal

import (
	"context"
	"sync"
)

type WorkerPool struct {
	WorkerCount  int
	JobStream    chan Job
	ResultStream chan Result
	Done         chan struct{}
}

func New(wc int) WorkerPool {
	return WorkerPool{
		WorkerCount:  wc,
		JobStream:    make(chan Job),
		ResultStream: make(chan Result),
		Done:         make(chan struct{}),
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.WorkerCount; i++ {
		wg.Add(1)

		go worker(ctx, &wg, wp.JobStream, wp.ResultStream)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.ResultStream)
}

func (wp WorkerPool) Results() <-chan Result {
	return wp.ResultStream
}

func (wp WorkerPool) GenerateJobBulk(jobBulk []Job) {
	for _, job := range jobBulk {
		wp.JobStream <- job
	}

	close(wp.JobStream)
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobStream <-chan Job, resultStream chan<- Result) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobStream:
			if !ok {
				return
			}

			resultStream <- job.execute(ctx)
		case <-ctx.Done():
			resultStream <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
