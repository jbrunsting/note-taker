package main

import (
	"fmt"
	"log"

	"github.com/jbrunsting/note-taker/manager"
	"github.com/jbrunsting/note-taker/request"
	"github.com/jbrunsting/note-taker/ui"
)

func main() {
	r := request.RequestFromArgs()
	m := manager.Manager{r.NotesDir}
	u := ui.UI{}

	if r.Cmd == request.NEW {
		if r.NotesDir == "" {
			log.Fatalf("TODO: error message, dir empty")
		}
		if r.NewArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		notes, err := m.ListNotes([]string{})
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
		// TODO: Should read notes to get highest ID
		header := fmt.Sprintf("[@%d", len(notes)+1)
		if len(r.Tags) != 0 {
			for _, tag := range r.Tags {
				header += ", #" + tag
			}
		}
		header += "]\n"
		err = m.CreateAndEdit(r.NewArgs.Title, header)
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
			notes, err := m.ListNotes(r.Tags)
			if err != nil {
				log.Fatalf("TODO: Error '%v'", err)
			}
			if len(notes) == 0 {
				log.Fatalf("TODO: No notes matched")
			}
			title = u.SearchForNote(notes)
			if title == "" {
				log.Fatalf("TODO: Title empty")
			}
		}

		err := m.Edit(title)
		if err != nil {
			log.Fatalf("Got error: '%v'", err)
		}
	} else if r.Cmd == request.BULK {
		if r.NotesDir == "" {
			log.Fatalf("TODO: error message, dir empty")
		}
		notes, err := m.ListNotes(r.Tags)
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
		manager.SortNotesById(notes)
		err = m.BulkEdit(notes)
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
	}
}
