package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fmt.Printf("What is your name? (2s)>> ")
	name, err := readLineTimeout(ctx, r)
	// Can't use `err == io.EOF` here as the error returned by `readLineTimeout`
	// is wrapped. This is a common idiom in Go, so one is always encouraged to
	// use `errors.Is` instead of "manually" comparing an error against another
	// using the `==` operator.
	if errors.Is(err, io.EOF) {
		fmt.Println()
		fmt.Println("got nothing")
		return
	}
	if err != nil {
		fmt.Println()
		log.Fatalf("failed to read input: %v", err)
	}
	fmt.Printf("got '%s'\n", name)
}

func readLineTimeout(ctx context.Context, r *bufio.Reader) (string, error) {
	type result struct {
		err error
		str string
	}
	// Buffered channel in order to avoid goroutine leak.
	c := make(chan result, 1)
	// Since `ReadString` is blocking, in order to make it "cancellable", one
	// needs to start a goroutine, making it run concurrently with respect to
	// the "cancellation check".
	go func() {
		s, err := r.ReadString('\n')
		if err != nil {
			c <- result{err: fmt.Errorf("failed to read: %w", err)}
			return
		}
		c <- result{str: strings.TrimSpace(s)}
	}()
	select {
	case res := <-c:
		return res.str, res.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
