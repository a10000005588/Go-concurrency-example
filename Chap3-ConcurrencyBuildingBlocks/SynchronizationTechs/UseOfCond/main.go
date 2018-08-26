package main

import (
	"sync"
)

func main() {

	// this cost CPU so much

	/*
		for condition() == false {}
	*/
	// this still unefficient

	/*  // 依然沒有效率，因為goroutine依然卡在OS thread中，只是多睡了1秒
	for conditinoTrue() == false {
		time.Sleep(1*time.Millisecond)
	}
	*/

	c := sync.NewCond(&sync.Mutex{})
	c.L.Lock()
	for conditionTrue() == false {
		// we wait to be notified that the condition has occured.

		// the call to wait doesn't just block, it "suspends" the current goroutine, allowing other goroutines to run on the OS thread.
		c.Wait()
	}
	c.L.Unlock()
}
