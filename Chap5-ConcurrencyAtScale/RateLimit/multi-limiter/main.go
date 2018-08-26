package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// 1. Here we define a "RateLimiter" interface so that a MultiLimiter can
//     recursively define other MultiLimter intances.
type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	// 2. We implement an optimization and sort by the Limit() of each RateLimiter
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

type multiLimiter struct {
	limiters []RateLimiter
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	// 3. Because we sort the "child RateLimiter" instances when multiLimiter
	//    is instantiated, we can simply return the most restrictive limit,
	//    which will be the "first element in the slice".
	return l.limiters[0].Limit()
}

func Open() *APIConnection {
	// Here we define our limit per second with no burstiness.
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
	// We define our limit per minute with a burstiness of 10 to give the users
	//   their initial pool. The limit per second will ensure we don't overload our system requests.
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)

	return &APIConnection{
		// We combine the two limits and set this as the master rate limiter for our APIConnection.
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

type APIConnection struct {
	rateLimiter RateLimiter
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
