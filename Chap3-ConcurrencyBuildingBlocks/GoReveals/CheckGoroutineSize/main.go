package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	memConsumed := func() uint64 {
		// create a runtime Garbage Collector.
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	// 1. We require a goroutine that will never exit so that we keep a number of them in a memory for measurement.
	noop := func() {
		wg.Done()
		<-c
	}

	const numGoroutines = 1e4 // 2. Here we define the number of goroutines to create.  Create  10000 goroutines.
	fmt.Println(numGoroutines)
	wg.Add(numGoroutines)
	before := memConsumed() // 3. We measure the amout of memory consumed before creating our goroutines.
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()                                         // 4. We measure the amount of memory consumed after creating our goroutines.
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000) // 1 kikabyte = 1000 bytes.

	// The results will show the size of gorouine is about 3 kb.
	// If a laptop has 8GB RAM.
	// You can create "millions of goroutines."
}
