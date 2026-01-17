package main

import (
	"fmt"
	"time"
)

func ping(pings chan<- string) {
	time.Sleep(1 * time.Second)
	pings <- "ping"
}

func pong(pings <-chan string, pongs chan<- string) {
	time.Sleep(1 * time.Second)
	select {
	case <-pings:
		pongs <- "pong"
	default:
		pongs <- "he doesn't pong me :("
	}

}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)

	go ping(pings)
	go pong(pings, pongs)

	fmt.Println(<-pongs)
}
