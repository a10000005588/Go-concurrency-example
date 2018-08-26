package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx := context.Background()
	// make client can have ability to signal to the context...
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	// after one second we cancel the request.
	defer cancel()
	// we pass value inside the ctx and pass it to server.
	//   which key is [foo], the value is [bar]
	ctx = context.WithValue(ctx, "foo", "bar")
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status)
	}
	io.Copy(os.Stdout, res.Body)
}
