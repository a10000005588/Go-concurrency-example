package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline.
// Fan-in is a term to describe the process of combining multiple results into one channel.

// key point to adopt the Fan-in,Fan-out pattern....
//  1. It doesn't rely on values that the stage had calculated before.
//  2. It takes a long time to run.

func main() {

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
	// convert the interace{} stream into integer stream
	toInt := func(
		done <-chan interface{},
		valueStream <-chan interface{},
	) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	primeFinder := func(
		done <-chan interface{},
		intStream <-chan int,
	) <-chan int {
		primeStream := make(chan int)
		go func() {
			defer close(primeStream)

			isPrime := func(value int) int {
				// In order to simulate that the stage is very computational, we just iterate the whole value to fine prime.
				// for i := 2; i <= int(math.Floor(math.Sqrt(float64(value)/2))); i++ {
				// for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
				for i := 2; i <= value; i++ {

					// If the value from random generator is not prime,
					//   then break this for loop and it necessary to return anything.
					if value%i == 0 {
						break
					}
				}
				return value
			}

			for rand := range intStream {
				select {
				case <-done:
					return
					// use isPrime to check whether the rand is prime or not.
					//   if is the prime number, then pass it to the primeStream
				case primeStream <- isPrime(rand):
				}
			}

		}()
		return primeStream
	}

	// take channel will limit the number of for loop to grab the value.
	take := func(
		done <-chan interface{},
		valueStream <-chan int,
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				// takeStream channel accept a channel which is only do output...
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	// Do fan-in...put all primeFinder into multiplex, combine them,
	//   and return a multiplexstream channel.
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan int,
		// 1. Use done channel to allow the goroutine to turn down.
		//      and the slice of interface{} channels to fan-in.
	) <-chan int {
		// 2. We create "sync.WaitGroup" so that we can wait until all channels have been drained.
		var wg sync.WaitGroup

		multiplexedStream := make(chan int)
		// 3. Here we create a function, "multiplex", which when passed a channel,
		//      will read from the channel, and pass the value read onto the multiplexedStream channel
		multiplex := func(c <-chan int) {
			// 每一個multiplex 執行完畢就會呼叫 wg.Done()
			defer wg.Done()
			// read the channel ... and put value into multiplexedStream channel
			// 讀取primeFinder channel，將從channel讀到的值，再丟回給multiplexedStream.
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// Select from all the channels
		// 4. This line increments the sync.WaitGroup by the number of channels we're multiplexing.
		wg.Add(len(channels))
		// with how many channels that create the corresponding number of multiplex.
		// 將primeFinder slice逐一讀出來，並啟動goroutine執行 multiplex(primeFinder)
		//   併發這些primeFinder
		//   只要這些primeFinder，全部的goroutine只要全部的goroutine塞入足夠take指定的數目，就會終止。
		for _, c := range channels {
			// put channel into multiplex...
			go multiplex(c)
		}

		// Wait for all the reads to complete
		// 5. Here we create a goroutine to wait for all the channels we're multiplexing to be
		//      drained so that we can close the multiplexedStream channel.
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	rand := func() interface{} { return rand.Intn(5000000000) }
	randIntStream := toInt(done, repeatFn(done, rand))

	// Do fan-out...
	// primeStream := primeFinder(done, randIntStream)

	// 有幾顆CPU、就copy幾份primeFinder stages
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders. \n", numFinders)
	// 創建一個finders slice來集合所有的primeFinder
	finders := make([]<-chan int, numFinders)
	fmt.Println("Primes:")

	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 40) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))

}
