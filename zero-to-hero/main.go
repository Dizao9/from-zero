package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	balance int
	mutex   sync.Mutex
}

func (bank *BankAccount) Deposit(wg *sync.WaitGroup, amount int) {
	defer wg.Done()
	bank.mutex.Lock()
	defer bank.mutex.Unlock()
	bank.balance += amount
}

// func worker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan int, slTime time.Duration) {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Halted by context")
// 			return
// 		case j, ok := <-jobs:
// 			if !ok {
// 				return
// 			}
// 			fmt.Println("worker", id, "start working with", j)
// 			select {
// 			case <-time.After(slTime * time.Second):
// 				fmt.Println("worker", id, "finish working with", j)
// 			case <-ctx.Done():

// 				return
// 			}
// 		}
// 	}
// }

func main() {
	var wg sync.WaitGroup
	bankAcc := &BankAccount{}

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go bankAcc.Deposit(&wg, 1)
	}

	// numJobs := 5
	// jobs := make(chan int, numJobs)

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// for w := 1; w <= 3; w++ {
	// 	wg.Add(1)
	// 	go worker(ctx, &wg, w, jobs, time.Duration(w))
	// }

	// for j := 1; j <= 3; j++ {
	// 	jobs <- j
	// }
	// fmt.Println("all jobs were sent")
	// close(jobs)

	wg.Wait()
	fmt.Println("Current balance = ", bankAcc.balance)
}
