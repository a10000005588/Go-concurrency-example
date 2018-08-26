package main

import (
	"sync"
	"testing"
)

// we start benchmark with the instruction:  "go test --bench=. --cpu=1"
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin // 1. we wait until we're told to begin
		for i := 0; i < b.N; i++ {
			c <- token // 2. we send the message to the receiver goroutines. A struct{}{} is called an empty struct and takes up no memory.
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin // 1.we wait until we're told to begin
		for i := 0; i < b.N; i++ {
			<-c // 3. Here we begin the performance timer.
		}
	}

	wg.Add(2)
	// 5. We tell the two goroutines to begin.
	go sender()
	go receiver()
	b.StartTimer() // 4. Here we begin the performance timer.
	close(begin)
	wg.Wait()

	// 225 ns per context switch. That's 0.225 us, or 92% faster than an OS context on the laptop machine.
}
