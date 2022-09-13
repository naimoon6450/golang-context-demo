package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	ctxClient := context.Background()
	// ctxClient, cancel := context.WithTimeout(ctxClient, time.Second)
	// defer cancel()
	// context with value

	req, err := http.NewRequest("GET", "http://localhost:8000", nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctxClient)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
