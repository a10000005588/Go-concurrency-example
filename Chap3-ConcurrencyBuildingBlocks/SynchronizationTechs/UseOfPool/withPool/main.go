package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func warmServiceConnCache() *sync.Pool {
	// Create a pool
	p := &sync.Pool{
		New: connectToService,
	}
	// put 10 connection thread to the pool
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a goroutine for handling the dialing th
	go func() {
		// get the sync.Pool
		connPool := warmServiceConnCache()

		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()

		wg.Done()

		for {

			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accpet connection: %v", err)
				continue
			}
			// get a connection from Pool
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			// if closing the connection, put it back to the pool.
			connPool.Put(svcConn)
			// then close.
			conn.Close()
		}
	}()
	return &wg
}

// benchmark with:
// go test --benchtime=10s --bench=.
