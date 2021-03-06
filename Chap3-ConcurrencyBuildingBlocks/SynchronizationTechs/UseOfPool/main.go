package main

import (
	"fmt"
	"sync"
)

// func main() {
// 	myPool := &sync.Pool{
// 		New: func() interface{} {
// 			fmt.Println("Creating new instance")
// 			return struct{}{}
// 		},
// 	}

// 	myPool.Get()
// 	instance := myPool.Get()
// 	myPool.Put(instance)
// 	myPool.Get()
// }

func main() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}
	// we only allocate 4 * 1024 = 4 KB for creating the
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte) // We are asserting the type is a pointer to a slice of bytes.
			defer calcPool.Put(mem)

			// Assume something interesting, but quick is being done with this memory.
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
