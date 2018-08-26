package log

import (
	"context"
	"log"
	"math/rand"
	"net/http"
)

type key int

// The key is a type that only the log pacakge can use !
const requestIDKey = key(42)

func Println(ctx context.Context, msg string) {
	id, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		log.Println("could not find request ID in context")
		return
	}
	log.Printf("[%d] %s", id, msg)
}

func Decorate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the context from request
		ctx := r.Context()
		// generate random number
		id := rand.Int63()
		// generate new context with value : requestIDKey
		ctx = context.WithValue(ctx, requestIDKey, id)
		// send new context back to server handler
		f(w, r.WithContext(ctx))
	}
}
