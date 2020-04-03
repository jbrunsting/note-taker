package main

import (
	"log"

	"github.com/jbrunsting/note-taker/manager"
	"github.com/jbrunsting/note-taker/request"
	"github.com/jbrunsting/note-taker/ui"
)

func main() {
	r := request.RequestFromArgs()
	m := manager.Manager{r.NotesDir}
    u := ui.UI{m}

	if r.Cmd == request.NEW {
		if r.NotesDir == "" {
			log.Fatalf("TODO: error message, dir empty")
		}
		if r.NewArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		header := ""
		if len(r.NewArgs.Tags) != 0 {
			header += "["
			for i, tag := range r.NewArgs.Tags {
				if i != 0 {
					header += ", "
				}
				header += "#" + tag
			}
			header += "]\n\n"
		}
		err := m.CreateAndEdit(r.NewArgs.Title, header)
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

		// TODO: Handle tags
		title := r.EditArgs.Title
		if title == "" {
			title = u.SearchForNote()
			if title == "" {
				log.Fatalf("TODO: Title empty")
			}
		}

		err := m.Edit(title)
		if err != nil {
			log.Fatalf("Got error: '%v'", err)
		}
	}
}
