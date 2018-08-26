package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		// 1. Here we set up a channel to send heartbeats on. We return this out of doWork.
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			defer close(heartbeat)
			defer close(results)
			// 2. Here we set the heartbeat to pulse at the pulseInterval we were given.
			//   Every pulseInterval there will be something to read on this channel.
			pulse := time.Tick(pulseInterval)
			// 3. This is just another ticker used to simulate work coming in.
			//   We choose a duration greater than the pulseInterval so that
			//   we can see some heartbeats coming out of the goroutine.
			//   模擬有工作排程的狀況
			workGen := time.Tick(2 * pulseInterval)
			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
					// 4. Not that we include a default clause. We must always guard against the fact that
					//   no one may be listening to our heartbeat. The results emitted from the goroutine
					//   are critical, but the pulses are not.
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
						// 5. Just like with "done" channels, anytime you perform a send or receive,
						//   you also need to include a case ofr the heartbeat's pulse.
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
					// 5.
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan interface{})
	// When pass 10 seconds, then close the done channels.
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	// the time of heartbeat is half of the timeout.
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
