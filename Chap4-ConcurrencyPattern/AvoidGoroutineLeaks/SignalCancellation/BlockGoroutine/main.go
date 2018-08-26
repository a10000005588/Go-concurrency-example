package main

import (
	"fmt"
	"math/rand"
)
// 該範例展示一個被confinement起來的producer goroutine沒有被告知要關閉！
func main() {
	// 創造一個 會建立各種新的亂數的channel ，並只有提供給外部scope的goroutine做讀取的動作
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		// 開始concurrency，自己獨立運作
		go func() {
			// 1. we print out the message when the goroutine successfully terminates.
			// 注意！ 在該範例中goroutine沒辦法順利被關閉，故該行不會被執行到
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()

		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	// 從randStream中拿出3個隨機值
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}
