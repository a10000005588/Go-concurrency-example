package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type startGoroutineFn func(
	done <-chan interface{},
	pulseInterval time.Duration,
) (heartbeat <-chan interface{}) // 1.

func main() {
	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn, // startGOroutine start the goroutine it's monitoring
	) startGoroutineFn { // The Steward itself returns a startGoroutineFB indicating that the steward itself is also monitorable.
		return func(
			done <-chan interface{},
			pulseInterval time.Duration,
		) <-chan interface{} {
			heartbeat := make(chan interface{})
			go func() {
				defer close(heartbeat)

				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}
				startWard := func() { // 3. we define a closure that encodes a consistent way to start the goroutine we're monitoring.
					wardDone = make(chan interface{})                   // 4. this is where we create a new channel that we'll pass into the ward goroutine in case we need to signal that it should halt.
					wardHeartbeat = startGoroutine(wardDone, timeout/2) // 5. Here we start the goroutine we'll be monitoring.
				}
				startWard()
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for {
					timeoutSignal := time.After(timeout)

					for { // 6. ensures that the steward can send out pulses of its own.
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case <-wardHeartbeat: // 7. If we receive the ward's pulse, we continue our monitoring loop.
							continue monitorLoop
						case <-timeoutSignal: // 8. indicates that if we don't receive a pulse from the ward within our time-out period,
							//    we request that the ward halt and we begin a new ward goroutine. We then continue monitoring.
							log.Println("steward: ward unhealthy; restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()
			return heartbeat
		}
	}

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

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWorkFn := func(
		done <-chan interface{},
		initList ...int,
	) (startGoroutineFn, <-chan interface{}) {
		intChanStream := make(chan (<-chan interface{}))
		intStream := bridge(done, intChanStream)
		doWork := func(
			done <-chan interface{},
			pulseInterval time.Duration,
		) <-chan interface{} {
			intStream := make(chan interface{})
			heartbeat := make(chan interface{})
			go func() {
				defer close(intStream)
				select {
				case intChanStream <- intStream:
				case <-done:
					return
				}

				pulse := time.Tick(pulseInterval)

				for {
				valueLoop:
					for _, intVal := range initList {
						/* 如果要接收多個，可以用以下寫法，不過要注意若剩下最後一個值，會發生 out of index
						for {
							intVal := initList[0]
							initList = initList[1:]
						*/
						if intVal < 0 {
							log.Printf("negative value: %v\n", intVal)
							return
						}

						for {
							select {
							case <-pulse:
								select {
								case heartbeat <- struct{}{}:
								default:
								}
							case intStream <- intVal:
								continue valueLoop
							case <-done:
								return
							}
						}
					}
				}
			}()
			return heartbeat
		}
		return doWork, intStream
	}

	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)
	doWorkWithSteward := newSteward(1*time.Millisecond, doWork)
	doWorkWithSteward(done, 1*time.Hour)

	for intVal := range take(done, intStream, 6) {
		fmt.Printf("Received: %v\n", intVal)
	}
}
