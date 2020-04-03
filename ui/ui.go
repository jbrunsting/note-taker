package ui

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jbrunsting/note-taker/manager"
)

const (
	titleColumnSize  = 40
	colorBrightWhite = 15
	colorBrightBlack = 8
)

type UI struct {
	NoteManager manager.Manager
}

func getEntry(note manager.Note) string {
	titleOutput := []rune(note.Title)
	if len(titleOutput) > titleColumnSize {
		titleOutput = titleOutput[:titleColumnSize]
		for i := 0; i < 3; i++ {
			titleOutput[titleColumnSize-i-1] = '.'
		}
	}
	return fmt.Sprintf(
		fmt.Sprintf("%%-%d.%ds  %%s", titleColumnSize, titleColumnSize),
		string(titleOutput),
		note.ModTime.Format("2006/01/02 15:04:05"),
	)
}

func (u *UI) SearchForNote() string {
	notes, err := u.NoteManager.ListNotes()
	if err != nil {
		log.Fatalf("TODO: Error")
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.TextStyle = ui.NewStyle(ui.ColorWhite, ui.Color(colorBrightBlack))
	l.SelectedRowStyle = ui.NewStyle(colorBrightWhite, ui.ColorBlack, ui.ModifierBold)
	l.WrapText = false
	l.SetRect(-1, -1, 80, 16)
	l.Border = false

	l.Rows = []string{}
	searchKey := ""
	sortRows := func() {
		l.Rows = []string{}
		manager.SortNotes(notes, searchKey)
		for _, note := range notes {
			l.Rows = append(l.Rows, getEntry(note))
		}
	}
	sortRows()
	l.ScrollBottom()

	t := widgets.NewParagraph()
	t.SetRect(-2, 15, 80, 16)
	t.Title = "> "
	t.Border = false

	render := func() {
		ui.Render(l)
		ui.Render(t)
	}
	render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-j>", "<Down>":
			l.ScrollDown()
		case "<C-k>", "<Up>":
			l.ScrollUp()
		case "<C-c>":
			return ""
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "<Enter>":
			return notes[l.SelectedRow].Title
		default:
			if e.Type == ui.KeyboardEvent {
				updated := true
				if e.ID == "<Space>" {
					searchKey += " "
				} else if e.ID == "<Backspace>" {
					if len(searchKey) != 0 {
						searchKey = searchKey[:len(searchKey)-1]
					}
				} else if len(e.ID) == 1 {
					searchKey += e.ID
				} else {
					updated = false
				}
				t.Title = fmt.Sprintf("> %s", searchKey)
				if updated {
					sortRows()
				}
			}
		}

		render()
	}
	return ""
}
