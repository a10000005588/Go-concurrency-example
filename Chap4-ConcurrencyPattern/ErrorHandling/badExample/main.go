package main

import (
	"fmt"
	"net/http"
)

func main() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		response := make(chan *http.Response)
		go func() {
			defer close(response)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("An error happens")
					// We see the goroutine doing its best to signal that there's an error.
					// and error can be passed back.
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case response <- resp:
				}
			}
		}()
		return response
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"http://www.google.com", "https//badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Println("Response: %v\n", response.Status)
	}
}
