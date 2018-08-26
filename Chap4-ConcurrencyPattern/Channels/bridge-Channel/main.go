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
	// 將一個channel的channel 拆解成 一個單純的channel
	// 簡單來說 Bridge 幫你先 for loop過channel內的channel值，在將值塞入到 bridge自己定義的channel，
	//   再回傳給 user，直接用一個 for loop 就可以搞定 <-chan <-chan 的channel 了
	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		// 1. This channel will return all value from bridge
		//   valStream channel會回傳所有來自bridge的值
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			// 2. This loop is responsible for pulling channels off of "chanStream"
			//   and providing them to a nested loop for use.
			//   該loop 負責將所有chanStream的channel值撈出來，
			//   並且給nested loop使用
			for {
				// steam channel 接收 從某個channel傳來的channel
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				// This loop is responsible for reading values off the channel it has been
				//   given and repeating those values onto "valStream"
				// When the stream we're currently looping over is closed,
				//   we break out of the loop performing the reads from this channel,
				//   and continue with the next iteration of the loop,
				//   selecting channels to read from.
				// This provides us with an unbroken stream of values.

				// 透過 orDone 將 stream的巢狀channel給包裝起來掉
				// 在遍覽整個 orDone 把值塞給 valStream.
				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	// genVals會回傳一個會回傳interface{} channel的channel
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				// 把stream channel 傳到 chanStream channel內
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
