package main

import "fmt"

func main() {

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

	tee := func(
		done <-chan interface{},
		in <-chan interface{},
	) (<-chan interface{}, <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				// 1. We will want to use local versions of out1 and out2, so we shadow these variables.
				var out1, out2 = out1, out2
				// 2. We are going to use one "select statement" so that writes to out1 and out2 don't block
				//   each other. To ensure both are written to, we'll perform two iterations of the
				//   select statement: one for each outbound channel.
				//   透過select statement out1 和 out2就不會互相卡住，以及透過for 我們可以保證out1和out2的運作
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						// 3. Once we've written to a channel, we set its shadowed copy to nil
						//   so that futher writes will block and the other channel may continue.
						//   透過將out1 out2 設置成nil ，可以阻止下一次的寫入
						//   並且保證其他channel可以順利運作，不會被block住
						//   也就是用 nil channel 來block，避免goroutine leak
						//     (詳情請見: https://medium.com/justforfunc/why-are-there-nil-channels-in-go-9877cc0b2308 )
						out1 = nil
					case out2 <- val:
						// 3. 同上
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
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

	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
