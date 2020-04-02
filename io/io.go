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
	lenWeight       = 1.0
	lenPenaltyFloor = 10.0
	lenPenaltyCeil  = 100.0
	matchingWeight  = 5.0
	dateWeight      = 3.0 // Date diff in days
	datePenaltyCeil = 50.0
	SecondsInDay    = 86400.0
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
	matchingScore := float64(getMatching(note.Title, searchKey)) * matchingWeight

	lenPenalty := float64(len(note.Title)) * lenWeight
	if lenPenalty < lenPenaltyFloor {
		lenPenalty = 0.0
	} else if lenPenalty > lenPenaltyCeil {
		lenPenalty = lenPenaltyCeil
	}

	datePenalty := (curTime.Sub(note.ModTime).Seconds() / SecondsInDay) * dateWeight
	if datePenalty > datePenaltyCeil {
		datePenalty = datePenaltyCeil
	}

	return matchingScore - lenPenalty - datePenalty
}

type Note struct {
	Title   string
	ModTime time.Time
}

func sortNotes(notes []Note, searchKey string) {
	t := time.Now()
	sort.SliceStable(notes, func(i, j int) bool {
		return getScore(notes[i], searchKey, t) > getScore(notes[j], searchKey, t)
	})
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
			note := Note{}
			note.Title = n[:len(n)-3]
			note.ModTime = f.ModTime()
			notes = append(notes, note)
		}
	}

	l := widgets.NewList()
	l.Title = "> "
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 80, 13)
	l.Border = false

	l.Rows = []string{}
	searchKey := ""
	sortRows := func() {
		l.Rows = []string{}
		sortNotes(notes, searchKey)
		for _, note := range notes {
			l.Rows = append(l.Rows, note.Title)
		}
	}

	sortRows()
	ui.Render(l)

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
			return "" // TODO: Return selected row
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
				l.Title = fmt.Sprintf("> %s", searchKey)
				if updated {
					sortRows()
				}
			}
		}

		ui.Render(l)
	}
	return ""
}
