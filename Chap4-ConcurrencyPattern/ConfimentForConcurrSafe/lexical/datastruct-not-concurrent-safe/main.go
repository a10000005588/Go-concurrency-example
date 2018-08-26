package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	// printData doesn't close around the "data slice", it cannot access it.
	// instead, it needs to take in a slice of byte to operate on.
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	// 這我們限制了data slice的切割範圍(沒有衝突的部分)給各自的goroutine，故goroutines不需要做synchronization的動作
	// 1. Here we pass in a slice containing the first three bytes in the data structure.
	go printData(&wg, data[:3])
	// 2. Here we pass in a slice containing the last three bytes in the data structure.
	go printData(&wg, data[3:])

	wg.Wait()
}
