package main

import (
	"fmt"
	"net/http"
)

func main() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan Result { // We return a channel that throw the Result (contain Error info) to other goroutine.
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				// Here, we wrap the error message with the response.
				result := Result{
					Error:    err,
					Response: resp,
				}
				if err != nil {
					fmt.Println("An error happens")
					// We see the goroutine doing its best to signal that there's an error.
					// and error can be passed back.
					fmt.Println(err)
					//continue
				}

				select {
				case <-done:
					return
					// we pass the Result to the channel which is handling the error message.
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	// We define a error counter.
	errCount := 0
	urls := []string{"a", "http://www.google.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			// If exceed the error number...
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
