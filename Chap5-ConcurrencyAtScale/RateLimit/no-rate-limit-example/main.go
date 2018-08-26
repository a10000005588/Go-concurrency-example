package main

import (
	"context"
	"log"
	"os"
	"sync"
)

func Open() *APIConnection {
	return &APIConnection{}
}

type APIConnection struct{}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// Pretend we do work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
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
