package main

import (
	"fmt"
	"time"
)

func main() {

	// 1. we pass "done" channel to the doWork function
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			// 負責check是否有關閉訊號過來
			for {
				select {
				// in our example, the strings channel input is nil.
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				// 2. Check whether our "done" channel has been signaled. If it has, we return from the goroutine.
				case <-done:
					fmt.Println("done channel receive signal")
					// 一旦return，就會觸發上面的defer，關閉 terminate channel
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)
	// 3. Another goroutine that will cancel the goroutine spawned in doWork.
	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		// will send message to done channel
		close(done)
	}()
	// 4. This is where we join the goroutine spawned from "doWork" with the main goroutine.
	// 這裡會block住，直到terminate channel被關閉或是有其他gouroutine向其取值
	<-terminated
	fmt.Println("Done.")
}
