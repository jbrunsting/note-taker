package io

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func SearchForNote(dir string) string {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "> "
	l.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
		"[10] bar",
		"[11] baz",
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 80, 13)
	l.Border = false
	ui.Render(l)

	searchKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-j>":
			l.ScrollDown()
		case "<C-k>":
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
		}

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
		}
		ui.Render(l)
	}
	return ""
}
