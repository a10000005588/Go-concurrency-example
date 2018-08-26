package main

import (
	"context"
	"log"
	"os"
	"sync"
	// please go get below first...
	"golang.org/x/time/rate"
)

/*
// Limit defines the maximum frequency of some events.
// Limit is represented as number of events per second.
// A zero Limit allows no events.
type Limit float64

// NewLimiter returns a new Limiter that allows events up to rate r
// and permits bursts of at most b tokens.
func NewLimiter(r Limit, b int) *NewLimiter

func Every(interval time.Duration) Limit)

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration/time.Duration(eventCount))
}
// Wait is shorthand for WaitN(ctx, 1)
func (lim *Limiter) Wait(ctx context.Context)
// WaitN blocks until lim permits n events to happen
// It returns an error if n exceeds the Limiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context'Deadline.
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)
*/

func Open() *APIConnection {
	return &APIConnection{
		// 設定 每次connection 的頻率為 1秒，並且 token 數量最多為1（表示只能一個goroutine來request)
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// 2. We set a.rateLimiter.Wait(ctx) to wait on the rate limiter
	//   to have enough access tokens for us to complete our request.
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// 2. same as ReadFile API
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func main() {
	// When 20 goroutines finish their work, print "Done"
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	// Create a api connection.
	apiConnection := Open()

	var wg sync.WaitGroup
	wg.Add(20)

	// Create 10 goroutines. Each goroutine start "ReadFile".
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot Readfile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	// Create another 10 goroutines. Each goroutine start "ResolveAddress"
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	wg.Wait()
}
