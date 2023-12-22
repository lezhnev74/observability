package main

import (
	"github.com/lezhnev74/observability"
	"log"
)

func main() {
	url := "http://127.0.0.1:8080/status"
	stats, err := observability.PollFpmStatus(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", stats)
}
