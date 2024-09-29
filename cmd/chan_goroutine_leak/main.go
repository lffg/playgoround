package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type state struct {
	buffered   bool
	unbuffered bool
}

var s state
var mu sync.Mutex

func runTest(name string, c chan struct{}, sb *bool) {
	log(name, "started")
	<-time.After(200 * time.Millisecond)
	log(name, "will send")
	// If this is an unbuffered channel, since there will be no one listening,
	// this send will never finish, making the goroutine leak.
	c <- struct{}{}
	log(name, "sent!")

	mu.Lock()
	defer mu.Unlock()
	*sb = true
	log(name, "finished")
}

func main() {
	buffered := make(chan struct{}, 1)
	unbuffered := make(chan struct{})

	go runTest("buffered", buffered, &s.buffered)
	go runTest("unbuffered", unbuffered, &s.unbuffered)

	select {
	case <-time.After(100 * time.Millisecond):
		log("main", "first arm finished")
	case <-buffered:
		// Won't happen since the first arm will finish first.
		unreachable()
	case <-unbuffered:
		// Won't happen since the first arm will finish first.
		unreachable()
	}

	// Even if we wait more time the unbuffered goroutine will never finish.
	<-time.After(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if !s.buffered {
		unreachable()
	}
	log("main", "buffered finished?= true")
	if s.unbuffered {
		unreachable()
	}
	log("main", "unbuffered finished?= false")
}

func log(name string, msg string) {
	fmt.Printf("[%s] %s\n", name, msg)
}

func unreachable() {
	fmt.Println("reached unreachable code")
	os.Exit(1)
}
