package main

import (
	"github.com/lezhnev74/observability"
	"log"
)

func main() {
	addr := "127.0.0.1:9090"
	path := "/status"

	stats, err := observability.PollFpmPoolStatus(addr, path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", stats)
}
