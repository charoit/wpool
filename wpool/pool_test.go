package wpool

import (
	"context"
	"testing"
	"time"
)

const (
	jobsCount   = 10
	workerCount = 2
)

func TestWorkerPool(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.AddFrom(testJobs())

	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}

			id, found := r.Metadata["id"]
			if !found {
				t.Fatal("id not found")
			}

			val := r.Value.(int)
			if val != id.(int)*2 {
				t.Fatalf("wrong value %v; expected %v", val, id.(int)*2)
			}
		case <-wp.Done:
			return
		}
	}
}

func TestWorkerPool_TimeOut(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond*10)
	defer cancel()

	go wp.Run(ctx)

	for {
		select {
		case r := <-wp.Results():
			if r.Err != nil && r.Err != context.DeadlineExceeded {
				t.Fatalf("expected error: %v; got: %v", context.DeadlineExceeded, r.Err)
			}
		case <-wp.Done:
			return
		}
	}
}

func TestWorkerPool_Cancel(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())

	go wp.Run(ctx)
	cancel()

	for {
		select {
		case r := <-wp.Results():
			if r.Err != nil && r.Err != context.Canceled {
				t.Fatalf("expected error: %v; got: %v", context.Canceled, r.Err)
			}
		case <-wp.Done:
			return
		}
	}
}

func testJobs() []Job {
	jobs := make([]Job, jobsCount)
	for i := 0; i < jobsCount; i++ {
		jobs[i] = Job{
			Metadata: map[string]interface{}{
				"id": i,
			},
			ExecFn: execFn,
			Args:   i,
		}
	}
	return jobs
}
