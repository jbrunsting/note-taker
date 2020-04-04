package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jbrunsting/note-taker/manager"
	"github.com/jbrunsting/note-taker/request"
	"github.com/jbrunsting/note-taker/ui"
)

func getTitle(title string, tags []string) {
}

func main() {
	r := request.RequestFromArgs()
	m := manager.Manager{r.NotesDir}
	u := ui.UI{}

	if r.NotesDir == "" {
		log.Fatalf("TODO: error message, dir empty")
	}

	if r.Cmd == request.NEW {
		if r.NewArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		notes, err := m.ListNotes([]string{})
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
		// TODO: Should read notes to get highest ID
		header := fmt.Sprintf("[@%d", len(notes)+1)
		if len(r.NewArgs.Tags) != 0 {
			for _, tag := range r.NewArgs.Tags {
				header += ", #" + tag
			}
		}
		header += "]\n"
		err = m.CreateAndEdit(r.NewArgs.Title, header)
		if err != nil {
			log.Fatalf("Got error: '%v'", err)
		}
	} else if r.Cmd == request.EDIT {
		if r.EditArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		title := r.EditArgs.Title
		if title == "" {
			notes, err := m.ListNotes(r.EditArgs.Tags)
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
	} else if r.Cmd == request.DELETE {
		if r.DeleteArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		title := r.DeleteArgs.Title
		if title == "" {
			log.Fatalf("TODO: Title empty")
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(fmt.Sprintf("Are you sure you want to delete %s (y/n): ", title))
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("TODO: Error %v", err)
		}

		if text[0] == 'y' && text[1] == '\n' {
			err := m.Delete(title)
			if err != nil {
				log.Fatalf("Got error: '%v'", err)
			}
			fmt.Printf("Deleted %s\n", title)
		} else {
			fmt.Printf("Did not delete\n")
		}
	} else if r.Cmd == request.CONCAT {
		if r.ConcatArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		notes, err := m.ListNotes(r.ConcatArgs.Tags)
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
		manager.SortNotesById(notes)
		err = m.ViewAll(notes)
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
	}
}
