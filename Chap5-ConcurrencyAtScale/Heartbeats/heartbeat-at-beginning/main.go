package main

import (
	"fmt"
	"math/rand"
)

func main() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		// 1. we create the heartbeat channel with a buffer of one.
		//   This ensure that "there's always at least one pulse"
		//   sent out even if no one is listening in time for send to occur.
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			for i := 0; i < 10; i++ {
				// 2. We set up a separate "select" block for the hearbeat.
				// We don't want to include this in the same select block as the send on results
				// because if the receiver isn't ready for the result, they'll receive pulse instead.
				// 這裡和result 的 select區分開來的原因是因為：
				//  若receiver尚未準備收到result時，會不小心收到heartbeat的訊號(pulse)
				select {
				case heartbeatStream <- struct{}{}:
				// 3. Once again we guard against the fact that no one may be listening to our heartbeats.
				// 由於heartbeat channel有buffer 1的size，故沒有設置default的話receiver又會收到pulse訊號
				default:
				}

				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()
		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}
