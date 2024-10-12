package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		// Notice that we must increment the wg counter *before* entering the
		// goroutine, otherwise `wg.Wait` might finish (with 0) before the
		// increment takes place. This might happen if some goroutine starts
		// executing too late.
		wg.Add(1)
		go func() {
			fmt.Printf("[%d] started\n", i)
			ch <- i
			fmt.Printf("[%d] sent!\n", i)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		// Below we're using a `range` loop over a channel. It will only "break"
		// *after* the channel gets closed, so we must do it here.
		// However, if we close the channel too early, some interested party may
		// not be able to send a message over it. Hence the usage of a wait
		// group to work as a barrier before closing it.
		close(ch)
	}()

	for i := range ch {
		fmt.Printf("[main] got from %d\n", i)
	}
}
