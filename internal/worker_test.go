package internal

import (
	"context"
	"testing"
	"time"
)

const (
	totalWorkers = 10
	totalJobs    = 1000
)

func TestWorkerPool(t *testing.T) {
	wp := New(totalWorkers)

	ctx, cancel := context.WithCancel(context.TODO())

	defer cancel()

	go wp.GenerateJobBulk(mockJobs())

	go wp.Run(ctx)

	for {
		select {
		case result, ok := <-wp.Results():
			if !ok {
				continue
			}

			if result.Err != nil {
				t.Fatalf("error occured: %v", result.Err)
			}

			if result.Value.(string) != "Tere" {
				t.Fatalf("result is not expected: %v", result.Value)
			}

		case <-wp.Done:
			return
		}

	}
}

func TestWorkerPool_Cancel(t *testing.T) {
	wp := New(totalWorkers)

	ctx, cancel := context.WithCancel(context.TODO())

	cancel()

	go wp.Run(ctx)

	for {
		select {
		case result, ok := <-wp.Results():
			if !ok {
				continue
			}

			if result.Err != context.Canceled {
				t.Fatalf("error was not called by context canceled %v", result.Err)
			}

		case <-wp.Done:
			return
		}

	}
}

func TestWorkerPool_Timeout(t *testing.T) {
	wp := New(totalWorkers)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond)

	cancel()

	go wp.Run(ctx)

	for {
		select {
		case result, ok := <-wp.Results():
			if !ok {
				continue
			}

			if result.Err != context.DeadlineExceeded {
				t.Fatalf("error was not called by context deadline exceeded %v", result.Err)
			}

		case <-wp.Done:
			return
		}

	}
}

func mockJobs() []Job {
	jobs := make([]Job, totalJobs)
	for i := 0; i < totalJobs; i++ {
		jobs[i] = Job{
			ID: i,
			ExecFunc: func(ctx context.Context, args interface{}) (interface{}, error) {
				return "Tere", nil
			},
			Args: i,
		}
	}

	return jobs
}
