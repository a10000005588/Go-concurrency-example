package main

func main() {
	for val := range myChan {
		// Do something with val
	}

	/* When there are too much select case, the code structure will be complex...
	   loop:
	   	for {
	   		select {
	   		case <-done:
	   			break loop
	   		case maybeVal, ok := <-myChan:
	   			if ok == false {
	   				return // or maybe break from for.
	   			}
	   			// Do something with val
	   		}
	   	}
	   }
	*/

	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	for val := range orDone(done, myChan) {
		// Do something with val
	}
}
