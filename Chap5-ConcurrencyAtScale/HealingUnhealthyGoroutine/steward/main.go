package main

import (
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

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, I'm irresponsible!")
		go func() {
			<-done
			log.Println("ward: I am halting")
		}()
		return nil
	}

	doWorkWithSteward := newSteward(4*time.Second, doWork)

	done := make(chan interface{})
	// 經過9秒後 關閉所有channel
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) {
	}
	log.Println("Done")
}
