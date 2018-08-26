package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	// Create a context with a cancel.
	ctx, cancel := context.WithCancel(ctx)
	/*
		// You can use "cancel" above to cancel the goroutine below.
		// After one second, you call the cancel...
		time.AfterFunc(time.Second, cancel)
	*/
	// We could simply say "context.WithTimeout" is same as above.
	context.WithTimeout(ctx, time.Second)
	// need to defer a cancel...
	// Notice, we do not defer the cancel, the context will stop mySleepAndTalk instantly.
	defer cancel()

	// After 5 seconds , it should print "hello"
	// sleepAndTalk(ctx, 5*time.Second, "hello")

	mySleepAndTalk(ctx, 5*time.Second, "hello")
}

func sleepAndTalk(ctx context.Context, d time.Duration, msg string) {
	time.Sleep(d)
	fmt.Println(msg)
}

func mySleepAndTalk(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Printf(msg)
		// if context is done, we log something
	case <-ctx.Done():
		log.Print(ctx.Err())
	}
}
