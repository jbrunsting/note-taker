package main

import (
	"github.com/jbrunsting/note-taker/editor"
	"github.com/jbrunsting/note-taker/io"
	"github.com/jbrunsting/note-taker/request"
	"log"
)

func main() {
	r := request.CreateRequest()
	log.Printf("Created request %v\n", r)

	if r.Cmd == request.NEW {
		if r.NotesDir == "" {
			log.Fatalf("TODO: error message, dir empty")
		}
		if r.NewArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		err := editor.CreateAndEdit(r.NotesDir, r.NewArgs.Title)
		if err != nil {
			log.Fatalf("Got error: '%v'", err)
		}
	} else if r.Cmd == request.EDIT {
		path := io.SearchForNote(r.NotesDir)
		log.Fatalf("Path found is %v", path)
	}
}
