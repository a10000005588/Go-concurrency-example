package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 該範例展示如何解決一個被confinement起來的producer goroutine沒有被告知要關閉的問題
// 透過給 producer goroutine with a channel informing it to exit.

func main() {
	// 創造一個 會建立各種新的亂數的channel ，並只有提供給外部scope的goroutine做讀取的動作
	// 並塞入一個負責接收關閉訊號的done channel.
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		// 開始concurrency，自己獨立運作
		go func() {
			// 1. we print out the message when the goroutine successfully terminates.
			// 注意！ 在該範例中goroutine沒辦法順利被關閉，故該行不會被執行到
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			// 透過 for-select 來監聽是否需要關閉channel的訊號
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	// 從randStream中拿出3個隨機值
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	// Simulate ongoing work
	time.Sleep(1 * time.Second)
}
