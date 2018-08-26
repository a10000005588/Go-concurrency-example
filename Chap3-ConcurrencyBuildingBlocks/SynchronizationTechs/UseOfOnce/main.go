package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)

	// AnotherOnce = sync.once ; AnotherOnce.Do觸發後，其他的AnotherOnce.DO就不會觸發！
	var AnotherCount int
	AnotherIncrement := func() { AnotherCount++ }
	AnotherDecrement := func() { AnotherCount-- }

	var AnotherOnce sync.Once
	AnotherOnce.Do(AnotherIncrement)
	AnotherOnce.Do(AnotherDecrement)

	fmt.Printf("Another: %d\n", AnotherCount)
}
