package io

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	titleColumnSize    = 40
	maxNonMatchPenalty = -100.0
	maxMatchScore      = 250.0
	dateWeight         = 10.0 // Date diff in days
	datePenaltyCeil    = 100.0
	SecondsInDay       = 86400.0
	ColorBrightWhite   = 15
	ColorBrightBlack   = 8
)

func getMatching(entry, searchKey string) int {
	entry = strings.ToLower(entry)
	searchKey = strings.ToLower(searchKey)
	sIndex := 0
	matching := 0
	for _, c := range entry {
		if sIndex >= len(searchKey) {
			break
		}
		if c == []rune(searchKey)[sIndex] {
			matching += 1
			sIndex += 1
		}
	}
	return matching
}

func getScore(note Note, searchKey string, curTime time.Time) float64 {
	numMatching := getMatching(note.Title, searchKey)
	percentMatching := float64(numMatching) / float64(len(note.Title))
	matchingScore := percentMatching*maxMatchScore + (1.0-percentMatching)*maxNonMatchPenalty

	datePenalty := (curTime.Sub(note.ModTime).Seconds() / SecondsInDay) * dateWeight
	if datePenalty > datePenaltyCeil {
		datePenalty = datePenaltyCeil
	}

	return matchingScore - datePenalty
}

type Note struct {
	Title   string
	ModTime time.Time
}

func fullMatch(s string, prefix string) bool {
	s = strings.ToLower(s)
	prefix = strings.ToLower(prefix)
	if len(s) < len(prefix) {
		return false
	}
	for i, c := range prefix {
		if c != []rune(s)[i] {
			return false
		}
	}
	return true
}

func sortNotes(notes []Note, searchKey string) {
	t := time.Now()
	sort.SliceStable(notes, func(i, j int) bool {
		iMatch := fullMatch(notes[i].Title, searchKey)
		jMatch := fullMatch(notes[j].Title, searchKey)
		if iMatch && !jMatch {
			return false
		} else if !iMatch && jMatch {
			return true
		}
		return getScore(notes[i], searchKey, t) < getScore(notes[j], searchKey, t)
	})
}

func getEntry(note Note) string {
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

func SearchForNote(dir string) string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	notes := []Note{}
	for _, f := range files {
		n := f.Name()
		if len(n) >= 3 && n[len(n)-3:] == ".md" {
			notes = append(notes, Note{n[:len(n)-3], f.ModTime()})
		}
	}

	l := widgets.NewList()
	l.TextStyle = ui.NewStyle(ui.ColorWhite, ui.Color(ColorBrightBlack))
	l.SelectedRowStyle = ui.NewStyle(ColorBrightWhite, ui.ColorBlack, ui.ModifierBold)
	l.WrapText = false
	l.SetRect(-1, -1, 80, 16)
	l.Border = false

	l.Rows = []string{}
	searchKey := ""
	sortRows := func() {
		l.Rows = []string{}
		sortNotes(notes, searchKey)
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
