package main

import (
	"fmt"
	"github.com/jbrunsting/note-taker/editor"
	"github.com/jbrunsting/note-taker/request"
	"log"
	"time"
)

func noteFilePath(dir string) string {
	timeStr := time.Now().Format("2006-01-02-15-04-05")
	return fmt.Sprintf("%v/%v.md", dir, timeStr)
}

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

		header := fmt.Sprintf("# %v", r.NewArgs.Title)
		editor.CreateAndEdit(noteFilePath(r.NotesDir), header)
	}
}
