package main

import (
	"fmt"
	"log"

	"github.com/jbrunsting/note-taker/editor"
	"github.com/jbrunsting/note-taker/io"
	"github.com/jbrunsting/note-taker/request"
)

func main() {
	r := request.CreateRequest()

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
		if r.NotesDir == "" {
			log.Fatalf("TODO: error message, dir empty")
		}
		if r.EditArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		title := r.EditArgs.Title
        if title == "" {
            title = io.SearchForNote(r.NotesDir)
            if title == "" {
                log.Fatalf("TODO: Title empty")
            }
        }
		path := fmt.Sprintf("%s/%s.md", r.NotesDir, title)

		err := editor.Edit(path)
		if err != nil {
			log.Fatalf("Got error: '%v'", err)
		}
	}
}
