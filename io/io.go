package io

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	LEN_PENALTY     = -10
	MIN_LEN_PENALTY = -125
	MATCHING_WEIGHT = 50
)

func getMatching(entry, searchKey string) int {
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

func getScore(entry, searchKey string) int {
	lenPenalty := len(entry) * LEN_PENALTY
	if lenPenalty < MIN_LEN_PENALTY {
		lenPenalty = MIN_LEN_PENALTY
	}

	return MATCHING_WEIGHT*getMatching(entry, searchKey) + lenPenalty
}

func sortNames(rows []string, searchKey string) {
	sort.Slice(rows, func(i, j int) bool {
		return getScore(rows[i], searchKey) > getScore(rows[j], searchKey)
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

	l := widgets.NewList()
	l.Title = "> "
	l.Rows = []string{}
	for _, f := range files {
		n := f.Name()
		if len(n) >= 3 && n[len(n)-3:] == ".md" {
			l.Rows = append(l.Rows, n[:len(n)-3])
		}
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 80, 13)
	l.Border = false

	searchKey := ""
	sortNames(l.Rows, searchKey)
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
				if e.ID == "<Space>" {
					searchKey += " "
				} else if e.ID == "<Backspace>" {
					if len(searchKey) != 0 {
						searchKey = searchKey[:len(searchKey)-1]
					}
				} else if len(e.ID) == 1 {
					searchKey += e.ID
				}
				l.Title = fmt.Sprintf("> %s", searchKey)
				sortNames(l.Rows, searchKey)
			}
		}

		ui.Render(l)
	}
	return ""
}
