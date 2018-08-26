package main

import (
	"fmt"
	"math/rand"
)

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

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
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

	done := make(chan interface{})
	defer close(done)

	// Define a first-class function
	rand := func() interface{} { return rand.Int() }

	// Create a for loop to retrieve  10 times of rand in repeatFn channel
	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}
