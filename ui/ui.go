package ui

import (
	"fmt"
	"log"
	"sort"
	"unicode"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jbrunsting/note-taker/manager"
)

const (
	maxSearchRows    = 1000
	titleColumnSize  = 40
	colorBrightWhite = 15
	colorBrightBlack = 8
)

type UI struct {
	Manager *manager.Manager
}

type TextSearchRow struct {
	SecondarySort int
	NoteTitle     string
	LineNum       int
	LineText      string
}

// Returns the number of characters read before a match
func charsOccurInOrder(line string, chars string) int {
	if len(line) == 0 || len(chars) == 0 {
		return -1
	}
	i := 0
	for li, c := range line {
		if unicode.ToLower(rune(c)) == unicode.ToLower(rune(chars[i])) {
			i += 1
			if i >= len(chars) {
				return li
			}
		}
	}
	return -1
}

func (u *UI) SearchForText(notes []manager.Note) string {
	searchRows := []TextSearchRow{}
	getRows := func(searchKey string) []string {
		searchRows = []TextSearchRow{}
		manager.SortNotes(notes, searchKey)
		for _, note := range notes {
			lines, err := u.Manager.ReadNote(&note)
			if err != nil {
				continue
			}

			for line, text := range lines {
				result := charsOccurInOrder(text, searchKey)
				if result != -1 {
					searchRows = append(
						searchRows,
						TextSearchRow{-result, note.Title, line, text},
					)
				}
				if len(searchRows) > maxSearchRows {
					break
				}
			}
			if len(searchRows) > maxSearchRows {
				break
			}
		}

		sort.SliceStable(searchRows, func(i, j int) bool {
			return searchRows[i].SecondarySort > searchRows[j].SecondarySort
		})

		rowText := []string{}
		for _, r := range searchRows {
			rowText = append(rowText, r.LineText)
		}

		return rowText
	}

	getResult := func(index int) string {
		if index < len(searchRows) {
			return searchRows[index].NoteTitle
		}
		return ""
	}

	return u.search(getRows, getResult)
}

func (u *UI) SearchForNote(notes []manager.Note) string {
	getRows := func(searchKey string) []string {
		rows := []string{}
		manager.SortNotes(notes, searchKey)
		for _, note := range notes {
			titleOutput := []rune(note.Title)
			if len(titleOutput) > titleColumnSize {
				titleOutput = titleOutput[:titleColumnSize]
				for i := 0; i < 3; i++ {
					titleOutput[titleColumnSize-i-1] = '.'
				}
			}
			rows = append(rows, fmt.Sprintf(
				fmt.Sprintf("%%-%d.%ds  %%s", titleColumnSize, titleColumnSize),
				string(titleOutput),
				note.ModTime.Format("2006/01/02 15:04:05"),
			))
		}
		return rows
	}

	getResult := func(index int) string {
		return notes[index].Title
	}

	return u.search(getRows, getResult)
}

func (u *UI) search(getRows func(string) []string, getResult func(int) string) string {
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

	searchKey := ""
	l.Rows = getRows(searchKey)
	if len(l.Rows) == 0 {
		l.Rows = append(l.Rows, "")
	}
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
			return getResult(l.SelectedRow)
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
					l.Rows = getRows(searchKey)
					if len(l.Rows) == 0 {
						l.Rows = append(l.Rows, "")
					}
				}
			}
		}

		render()
	}
	return ""
}
