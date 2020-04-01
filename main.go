package main

import (
	"log"
	"github.com/jbrunsting/note-taker/request"
	"github.com/jbrunsting/note-taker/editor"
)

func main() {
    r := request.CreateRequest()
    log.Printf("Created request %v\n", r)

	if r.Cmd == request.NEW {
        editor.CreateAndEdit("/home/jacob/temp/test.txt")
	}
}
