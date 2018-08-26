package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	a := asChain(1, 3, 5, 7)
	b := asChain(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}

func asChain(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

// 1. Will casue infinite 0 return due to close channel

// func merge(a, b <-chan int) <-chan int {
// 	c := make(chan int)
// 	go func() {
// 		for {
// 			select {
// 			case v := <-a:
// 				c <- v
// 			case v := <-b:
// 				c <- v
// 			}
// 		}
// 	}()
// 	return c
// }

// 2. Will cause deadlock because we want retrieve something from nil channel [c] when there is no input to [c]

// func merge(a, b <-chan int) <-chan int {
// 	c := make(chan int)
// 	go func() {
// 		adone, bdone := false, false
// 		for !adone || !bdone {
// 			select {
// 			case v, ok := <-a:
// 				if !ok {
// 					adone = true
// 					continue
// 				}
// 				c <- v
// 			case v, ok := <-b:
// 				if !ok {
// 					bdone = true
// 					continue
// 				}
// 				c <- v
// 			}
// 		}
// 	}()
// 	return c
// }

// 3. Will cause looping the  case v,ok := <-a if the b channel is still not finished
// func merge(a, b <-chan int) <-chan int {
// 	c := make(chan int)
// 	go func() {
// 		defer close(c) // when finishing iterating [a] and [b] channel, close the [c] channel to avoid someone retrive nil from [c] casuing deadlock.
// 		adone, bdone := false, false
// 		for !adone || !bdone {
// 			select {
// 			case v, ok := <-a:
// 				if !ok {  // ok == false means the channel is close
// 					adone = true
// 					log.Printf("a is done")
// 					continue
// 				}
// 				c <- v
// 			case v, ok := <-b:
// 				if !ok {
// 					bdone = true
// 					log.Printf("b is done")
// 					continue
// 				}
// 				c <- v
// 			}
// 		}
// 	}()
// 	return c
// }

/*
func main() {
	var c chan int
	<-c  // if we retrieve the nil channel that will cause block!

	close(c)  // when we got someting from a close channel will cause panic.
}
*/

// This time , we don't want to check again for the finished a or b channel
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c) // when finishing iterating [a] and [b] channel, close the [c] channel to avoid someone retrive nil from [c] casuing deadlock.
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil // make a to be a nil channel
					log.Printf("a is done")
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil // make b to be a nil channel
					log.Printf("b is done")
					continue
				}
				c <- v
			}
		}
	}()
	return c
}
