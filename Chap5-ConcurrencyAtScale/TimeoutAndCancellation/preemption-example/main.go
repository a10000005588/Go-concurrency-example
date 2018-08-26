package main

func main() {
	var value interface{}
	select {
	case <-done:
		return
	case value = <-valueStream:
	}

	// we should make reallyLongCalculation be preemptive via putting done channel.
	reallyLongCalculation := func(
		done <-chan interface{},
		value interface{},
	) interface{} {
		// make the longCalculation also be preemptive too.
		intermediateResult := longCalculation(done, value)
		select {
		case <-done:
			return nil
		default:
		}

		return longCalculation(done, intermediateResult)
	}

	result := reallyLongCalculation(value)

	select {
	case <-done:
		return
	case resultStream <- result:
	}
}
