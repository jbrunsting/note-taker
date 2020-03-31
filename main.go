package main

import (
	"log"
	"github.com/jbrunsting/note-taker/request"
)

func main() {
    r := request.CreateRequest()
    log.Printf("Created request %v\n", r)
}
