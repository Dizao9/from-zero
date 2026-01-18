package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan int, slTime time.Duration) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Halted by context")
			return
		case j, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Println("worker", id, "start working with", j)
			select {
			case <-time.After(slTime * time.Second):
				fmt.Println("worker", id, "finish working with", j)
			case <-ctx.Done():

				return
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup
	numJobs := 5
	jobs := make(chan int, numJobs)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(ctx, &wg, w, jobs, time.Duration(w))
	}

	for j := 1; j <= 3; j++ {
		jobs <- j
	}
	fmt.Println("all jobs were sent")
	close(jobs)

	wg.Wait()
}
