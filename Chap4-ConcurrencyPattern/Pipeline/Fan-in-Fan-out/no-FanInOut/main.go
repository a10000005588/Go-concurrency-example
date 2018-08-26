package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline.
// Fan-in is a term to describe the process of combining multiple results into one channel.

// key point to adopt the Fan-in,Fan-out pattern....
// 1. It doesn't rely on values that the stage had calculated before.
// 2. It takes a long time to run.

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
					// then break this for loop and it necessary to return anything.
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
					// if is the prime number, then pass it to the primeStream
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
	) <-chan int {
		takeStream := make(chan int)

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

	rand := func() interface{} { return rand.Intn(5000000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 40) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}
