package main

func main() {
	var value interface{}
	select {
	case <-done:
		return
	case value = <-valueStream:
	}
	// it will take a long time...
	// and it is no-preemption.
	result := reallyLongCalculation(value)

	select {
	case <-done:
		return
	case resultStream <- result:
	}
}
