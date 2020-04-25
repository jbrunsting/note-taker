package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jbrunsting/note-taker/html"
	"github.com/jbrunsting/note-taker/manager"
	"github.com/jbrunsting/note-taker/request"
	"github.com/jbrunsting/note-taker/ui"
)

func defaultNotesDir() string {
	home := os.Getenv("HOME")
	path := home + "/.note-taker"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	return path
}

func main() {
	r := request.RequestFromArgs()

	if r.NotesDir == "" {
		r.NotesDir = defaultNotesDir()
	}

	m := manager.Manager{r.NotesDir}
	u := ui.UI{&m}

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
	} else if r.Cmd == request.MV {
		if r.MvArgs == nil {
			log.Fatalf("TODO: No image thing")
		}
		if r.MvArgs.Title == "" || r.MvArgs.Src == "" {
			log.Fatalf("TODO: Better arg validation")
		}
		components := strings.Split(r.MvArgs.Src, ".")
		m.Move(r.MvArgs.Src, r.MvArgs.Title, components[len(components)-1])
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
				fmt.Printf("No notes found\n")
                os.Exit(1)
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
	} else if r.Cmd == request.FIND {
		if r.FindArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		notes, err := m.ListNotes(r.FindArgs.Tags)
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
		title := u.SearchForText(notes)
		if title == "" {
			log.Fatalf("TODO: Title empty")
		}

		err = m.Edit(title)
		if err != nil {
			log.Fatalf("Got error: '%v'", err)
		}
	} else if r.Cmd == request.HTML {
		if r.HtmlArgs == nil {
			log.Fatalf("TODO: error message, shouldn't get here")
		}

		notes, err := m.ListNotes(r.HtmlArgs.Tags)
		if err != nil {
			log.Fatalf("TODO: Error '%v'", err)
		}
		manager.SortNotesById(notes)
		o, err := html.GenerateHTML(notes, r.NotesDir)
		if err != nil {
			// TODO: Add err check function that logs error nicely
			log.Fatalf("TODO: Error '%v'", err)
		}
		if r.HtmlArgs.File == "" {
			fmt.Println(string(o))
		} else {
			err := ioutil.WriteFile(r.HtmlArgs.File, []byte(o), 0644)
			if err != nil {
				log.Fatalf("TODO: Error '%v'", err)
			}
		}
	}
}
