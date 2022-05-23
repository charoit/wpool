package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"wpool/wpool"
)

const workerCount = 10

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wp := wpool.New(workerCount)

	go jobGenerator(ctx, &wp)
	go wp.Run(ctx)

loop:
	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}

			id, found := r.Metadata["id"]
			if !found {
				log.Fatal("id not found")
			}

			val := r.Value.(int)
			if val != id.(int)*2 {
				log.Fatalf("wrong value %v; expected %v", val, id.(int)*2)
			}
			log.Println(r)

		case <-wp.Done:
			break loop

		case <-ctx.Done():
			stop()
			break loop
		}
	}

	fmt.Println("Done!")
}

func jobGenerator(ctx context.Context, wp *wpool.WorkerPool) {
	count := 1

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		default:
			wp.Add(wpool.Job{
				Args: count,
				ExecFn: func(ctx context.Context, args interface{}) (interface{}, error) {
					argVal, ok := args.(int)
					if !ok {
						return nil, errors.New("wrong argument type")
					}

					return argVal * 2, nil
				},
				Metadata: map[string]interface{}{
					"id": count,
				},
			})

			count++
		}
	}
}
