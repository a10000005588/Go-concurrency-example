package main

import "fmt"

// hocs : means higher order componet which means a function can accept a component and also return a component.
func main() {
	data := make([]int, 4)

	// loopData can access data[]
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	// the loop of handleData can also access data[]
	for num := range handleData {
		fmt.Println(num)
	}
}
