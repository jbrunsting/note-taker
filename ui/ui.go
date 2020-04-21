package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"syscall"
	"unicode"

	"github.com/jbrunsting/note-taker/manager"
	"golang.org/x/sys/unix"
)

const (
	maxSearchRows    = 1000
	titleColumnSize  = 40
	colorBrightWhite = 15
	colorBrightBlack = 8
	enter            = 10
	del              = 127
	altB0            = 27
	altB1            = 91
	downArrowB2      = 66
	upArrowB2        = 65
	shiftTabB2       = 90
	tab              = 9
	minPrintable     = 32
	maxPrintable     = 126
	rowsToShow       = 15
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

func min(i int, j int) int {
	if i < j {
		return i
	}
	return j
}

// Returns the number of rows printed
func printSearch(rows []string, selectedRow int, searchKey string) int {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		log.Fatalf("TOOD: Err %v", err)
	}
	maxWidth := int(ws.Col)

	fmt.Printf("\n")
	rowsPrinted := 1

	// Print blank lines so we fill 15 lines even with less results
	for i := 0; i < rowsToShow-len(rows); i++ {
		fmt.Printf("\n")
		rowsPrinted += 1
	}

	// Print in reverse order so the best result is at the bottom
	for i := min(len(rows), rowsToShow) - 1; i >= 0; i-- {
		if i == selectedRow {
			fmt.Printf(">")
		} else {
			fmt.Printf(" ")
		}
		row := rows[i]
		if len(row) >= maxWidth-5 {
			row = row[:maxWidth-8]
			for i := 0; i < 3; i++ {
				row += "."
			}
		}
		fmt.Printf(" %s\n", row)
		rowsPrinted += 1
	}

	fmt.Printf("> %s", searchKey)

	return rowsPrinted
}

func (u *UI) search(getRows func(string) []string, getResult func(int) string) string {
	rows := []string{}
	searchKey := ""
	selectedRow := 0
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		os.Exit(1)
	}()

	prevRowsPrinted := 0
	var b []byte = make([]byte, 3)
	for {
		rows = getRows(searchKey)
		if selectedRow >= len(rows) {
			selectedRow = len(rows) - 1
			if selectedRow < 0 {
				selectedRow = 0
			}
		}

		fmt.Printf("\r\033[K")
		for i := 0; i < prevRowsPrinted; i++ {
			fmt.Printf("\033[1A\033[K")
		}
		prevRowsPrinted = printSearch(rows, selectedRow, searchKey)

		os.Stdin.Read(b)
		if b[0] == enter {
			break
		} else if b[0] == del {
			if len(searchKey) > 0 {
				searchKey = searchKey[:len(searchKey)-1]
			}
		} else if minPrintable <= int(b[0]) && int(b[0]) <= maxPrintable {
			searchKey += string(b[0])
		} else if len(b) == 3 && b[0] == altB0 && b[1] == altB1 && (b[2] == upArrowB2 || b[2] == shiftTabB2) {
			// Reverse direction because UI is bottom up
			selectedRow += 1
			if selectedRow >= rowsToShow {
				selectedRow = rowsToShow - 1
			}
		} else if (len(b) == 3 && b[0] == altB0 && b[1] == altB1 && b[2] == downArrowB2) || b[0] == tab {
			selectedRow -= 1
			if selectedRow < 0 {
				selectedRow = 0
			}
		}
	}

	if selectedRow < 0 || selectedRow >= len(rows) {
		return ""
	}
	return getResult(selectedRow)
}
