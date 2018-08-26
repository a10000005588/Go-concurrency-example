package main

import (
	"fmt"
	"time"
)

// 該範例透過or-channel來統一集中所有要關閉其他channel的 "done" channels
func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	// 1. 傳入 slices of channels，回傳單一個channel
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {

		// 2. 既然為一個recursive function，需要設置中斷條件
		// 所以如果slice lenght = 0 那就停止
		case 0:
			return nil
		// 3. 如果slice只有包含一個channel，就直接回傳
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		// 4. 我們該範例主要的邏輯，為recursive的地方
		//  where we can wait for messages on our channels withouth "blocking"
		go func() {
			defer close(orDone)
			switch len(channels) {
			// 5. Because we're recursing, every recursive call to "or" will at least have two channels.
			// As an optimization to keep the number of goroutines constrained,
			//   we place a special case here for calls to "or" with only two channels.
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			// 6. We recursively create an "or-channel" from all the channels in our slice after the third index,
			//   and third index, and then select from this.
			//  在channels[n] 第三個index之後塞入給or()做recursively.
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
					// 再次呼叫or func()...做recursives...
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
	// 1. This function simply creates a channel that will close when the time specified in the after elapses.
	// 建立一個會在設定的時間過後就關閉的channel
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			// 在時間結束後關閉 channel c
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	// 2. We keep track of roughly when the channel from the "or" function begins to block.
	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	// 3. We print the time it took for the read to occur.
	fmt.Printf("done after %v", time.Since(start))

}
