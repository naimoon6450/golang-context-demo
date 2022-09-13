package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

type username string

const userKey username = "username"

func main() {
	http.HandleFunc("/", handler)
	// print listening on port 8080
	log.Printf("listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// ideally the username comes from the request object, but hardcoding here
	ctx := context.WithValue(r.Context(), userKey, "nsiraj")
	log.Printf("handler started")
	defer log.Printf("handler ended")

	// print context value
	log.Printf("grettings from handler: %v", ctx.Value(userKey))
	serviceA(ctx, w)

	log.Printf("number of go routines: %d", runtime.NumGoroutine())
}

func serviceA(ctx context.Context, w http.ResponseWriter) {
	// do something and pass it to serviceB
	log.Printf("grettings from svc A: %v", ctx.Value(userKey))
	for {
		select {
		case <-time.After(1 * time.Second):
			log.Print("finished processing in serviceA")
			serviceBContextAware(ctx, w)
			// serviceBNonContextAware(w)
			return
		case <-time.After(5 * time.Second): // server level request limits
			// write to writer
			log.Print("request timeout serviceA")
			return
		case <-ctx.Done():
			ctxErr := ctx.Err()
			log.Printf("serviceA: %v", ctxErr)
			return
		}
	}
}

func serviceBNonContextAware(w http.ResponseWriter) {
	count := 0
	timeTick := time.Tick(1 * time.Second)
	timeAfter := time.After(30 * time.Second)
	for {
		select {
		case <-timeAfter:
			log.Print("request timeout serviceB")
			return
		case <-timeTick:
			log.Print("processing in serviceB")
			count++
			// print count
			log.Printf("count: %d", count)
			if count == 5 {
				// write to writer
				log.Print("finished processing in serviceB")
				serviceCNonaware()
				return
			}
		}
		log.Print("still processing in serviceB")
		time.Sleep(1 * time.Second)
	}
}

func serviceBContextAware(ctx context.Context, w http.ResponseWriter) {
	log.Printf("grettings from svc B: %v", ctx.Value(userKey))
	count := 0
	timeTick := time.Tick(1 * time.Second)
	timeOut := time.After(5 * time.Second)
	for {
		select {
		case <-timeTick:
			log.Print("processing in serviceB")
			count++
			// print count
			log.Printf("count: %d", count)
			if count == 5 {
				serviceCAware(ctx, w)
				return
			}
		case <-timeOut:
			log.Print("request timeout serviceB")
			return
		case <-ctx.Done():
			ctxErr := ctx.Err()
			log.Printf("serviceB: %v", ctxErr)
			return
		}
		log.Print("still processing in serviceB")
		// time.Sleep(5 * time.Second)
	}
}

func serviceCNonaware() {
	// do something
	log.Print("start processing in serviceC")
	time.Sleep(5 * time.Second)
	log.Print("finished processing in serviceC")
}

func serviceCAware(ctx context.Context, w http.ResponseWriter) {
	// do something
	log.Printf("grettings from svc C: %v", ctx.Value(userKey))
	log.Print("start processing in serviceC")
	time.Sleep(5 * time.Second)

	log.Print("finished processing in serviceC")
	// write to writer
	fmt.Fprintln(w, "hello from C")

	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		log.Printf("serviceB: %v", ctxErr)
		return
	case <-time.After(5 * time.Second):
		log.Print("request timeout serviceB")
		return
	}
}
