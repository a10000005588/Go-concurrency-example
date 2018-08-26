package main

import "fmt"

func main() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(completed)
			// because the strings channel is nil, will never actually gets any strings written onto it.
			// and the goroutine containing "doWork" will remain in memory for the lifetime of this process.
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
		}()
		return completed
	}
	// throw a nil to channel
	doWork(nil)
	fmt.Println("Done.")
}
