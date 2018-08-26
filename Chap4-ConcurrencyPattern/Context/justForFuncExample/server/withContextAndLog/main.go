package main

import (
	"context"
	"fmt"
	// "log"
	"net/http"
	"time"
	// import the self-defined log package...
	"Gopher/Chap4-GoConcurrencyPattern/Context/justForFuncExample/log"
)

func main() {
	http.HandleFunc("/", log.Decorate(handler))
	panic(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// get the ctx from client
	ctx := r.Context()
	// will contain same value 42, with value 100 all the time.
	ctx = context.WithValue(ctx, int(42), int64(100))

	// log.Println will append id to message.
	log.Println(ctx, "handler started")
	defer log.Println(ctx, "handler ended")

	// handle the ctx value from client which key is [foo]
	fmt.Printf("value for foo is %v", ctx.Value("foo"))

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "hello")
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
