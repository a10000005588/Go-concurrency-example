package main

import "fmt"

func main() {
	chanOwner := func() <-chan int {
		// 1. we intantiate the channel within the lexical scope of the "chanOwner"
		//    this limit the scope of the wirte aspect of the "results" channel
		// (in other words) it 'confines' the write aspect of this channel to prevent other goroutines from writing to it.
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	// 3. We receive a read-only copy of an int channel.
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	// 2. We receive the read aspect of the channel and we're able to pass it into the consumer.
	results := chanOwner()
	consumer(results)
}
