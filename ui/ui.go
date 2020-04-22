package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"unicode"

	"github.com/jbrunsting/note-taker/manager"
	"golang.org/x/sys/unix"
)

const (
	maxSearchRows   = 1000
	titleColumnSize = 25
	enter           = 10
	del             = 127
	altB0           = 27
	altB1           = 91
	downArrowB2     = 66
	upArrowB2       = 65
	shiftTabB2      = 90
	tab             = 9
	minPrintable    = 32
	maxPrintable    = 126
	rowsToShow      = 15
)

type UI struct {
	Manager *manager.Manager
}

type textSearchRow struct {
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

func getSearchTextRowComponents(text string, chars string) []RowComponent {
	components := []RowComponent{}

	curMatch := ""
	curNonMatch := ""

	i := 0
	for ti, c := range text {
		if unicode.ToLower(rune(c)) == unicode.ToLower(rune(chars[i])) {
			if curNonMatch != "" {
				components = append(components, RowComponent{curNonMatch, RowText, -1, -1})
				curNonMatch = ""
			}
			curMatch += string(c)
			i += 1
			if i >= len(chars) {
				if ti+1 < len(text) {
					curNonMatch = text[ti+1:]
				}
				break
			}
		} else {
			if curMatch != "" {
				components = append(components, RowComponent{curMatch, RowSelectedText, -1, -1})
				curMatch = ""
			}
			curNonMatch += string(c)
		}
	}

	if curMatch != "" {
		components = append(components, RowComponent{curMatch, RowSelectedText, -1, -1})
	}
	if curNonMatch != "" {
		components = append(components, RowComponent{curNonMatch, RowText, -1, -1})
	}

	// Replace the start with an elipse if it is not a match
	if components[0].Type == RowText {
		components[0].Text = "..."
	}

	return components
}

func (u *UI) SearchForText(notes []manager.Note) string {
	searchRows := []textSearchRow{}
	getRows := func(searchKey string) [][]RowComponent {
		searchRows = []textSearchRow{}
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
						textSearchRow{-result, note.Title, line, text},
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

		rows := make([][]RowComponent, 0)
		for _, r := range searchRows {
			rowComponents := []RowComponent{}
			rowComponents = append(rowComponents, RowComponent{r.NoteTitle, RowTitle, titleColumnSize, titleColumnSize})
			rowComponents = append(rowComponents, RowComponent{"", RowDecoration, 1, 1})
			rowComponents = append(rowComponents, getSearchTextRowComponents(r.LineText, searchKey)...)
			rows = append(rows, rowComponents)
		}

		return rows
	}

	getResult := func(index int) string {
		if index < len(searchRows) {
			return searchRows[index].NoteTitle
		}
		return ""
	}

	return u.SearchList(getRows, getResult)
}

func (u *UI) SearchForNote(notes []manager.Note) string {
	getRows := func(searchKey string) [][]RowComponent {
		manager.SortNotes(notes, searchKey)

		rows := make([][]RowComponent, 0)
		for _, note := range notes {
			rowComponents := []RowComponent{}
			rowComponents = append(rowComponents, RowComponent{note.Title, RowTitle, titleColumnSize, titleColumnSize})
			rowComponents = append(rowComponents, RowComponent{"", RowDecoration, 1, 1})
			rowComponents = append(rowComponents, RowComponent{note.ModTime.Format("2006/01/02 15:04:05"), RowDate, -1, -1})
			rows = append(rows, rowComponents)
		}

		return rows
	}

	getResult := func(index int) string {
		return notes[index].Title
	}

	return u.SearchList(getRows, getResult)
}

func min(i int, j int) int {
	if i < j {
		return i
	}
	return j
}

const (
	RowTitle = iota
	RowText
	RowSelectedText
	RowDate
	RowDecoration
)

type RowComponent struct {
	Text     string
	Type     int
	MinWidth int
	MaxWidth int
}

func willPrint(text string, pos int) bool {
	if !(minPrintable <= int(text[pos]) && int(text[pos]) <= maxPrintable) {
		return false
	}
	// Check if it is a color code
	prevEscPos := -1
	for i := pos; i >= 0; i-- {
		if int(text[i]) == 27 {
			prevEscPos = i
			break
		}
	}
	if prevEscPos != -1 {
		for i := prevEscPos; i < len(text); i++ {
			if text[i] == 'm' {
				if i >= pos {
					return false
				}
				break
			}
		}
	}

	return true
}

func constrainText(text string, minWidth int, maxWidth int, fgcolor string, bgcolor string, elipsize bool) string {
	width := 0
	output := ""
	if fgcolor != "" {
		output += fmt.Sprintf("\u001b[%sm", fgcolor)
	}
	if bgcolor != "" {
		output += fmt.Sprintf("\u001b[%sm", bgcolor)
	}
	for ti, c := range text {
		if maxWidth != -1 && willPrint(text, ti) {
			if elipsize && width+3 == maxWidth {
				output += "..."
				width = maxWidth
			} else if width < maxWidth {
				output += string(c)
				width += 1
			}
		} else {
			output += string(c)
		}
	}

	if minWidth != -1 {
		for i := 0; i < minWidth-width; i++ {
			output += " "
		}
	}

	return output
}

// Returns the number of rows printed
func printSearch(rows [][]RowComponent, selectedRow int, searchKey string) int {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		log.Fatalf("TOOD: Err %v", err)
	}
	screenWidth := int(ws.Col)

	rowsPrinted := 0

	// Print blank lines so we fill 15 lines even with less results
	for i := 0; i < rowsToShow-len(rows); i++ {
		fmt.Printf("\n")
		rowsPrinted += 1
	}

	// Print in reverse order so the best result is at the bottom
	topRow := min(len(rows)-1, rowsToShow-1)
	if topRow < selectedRow {
		topRow = min(len(rows)-1, selectedRow)
	}
	for i := topRow; i >= topRow-rowsToShow+1 && i >= 0; i-- {
		line := ""
		if i == selectedRow {
			line += "> "
		} else {
			line += "  "
		}

		for _, component := range rows[i] {
			fgcolor := "37"
			if component.Type == RowSelectedText {
				fgcolor = "91"
			}
			if i == selectedRow {
				fgcolor += ";1"
			}
			line += constrainText(component.Text, component.MinWidth, component.MaxWidth, fgcolor, "", true)
		}

		bgcolor := ""
		if i == selectedRow {
			bgcolor = "40"
		} else {
			bgcolor = "100"
		}
		fmt.Printf("%s\n", constrainText(line, 0, screenWidth, "", bgcolor, true))
		rowsPrinted += 1
	}

	fmt.Printf("> %s", searchKey)

	return rowsPrinted
}

func (u *UI) SearchList(getRows func(string) [][]RowComponent, getResult func(int) string) string {
	var rows [][]RowComponent
	searchKey := ""
	selectedRow := 0
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	prevRowsPrinted := 0
	var b []byte = make([]byte, 3)
	for {
		rows = getRows(searchKey)
		if selectedRow >= len(rows) {
			selectedRow = len(rows) - 1
		} else if selectedRow < 0 {
			selectedRow = 0
		}

		fmt.Printf("\u001b[0m\r\033[K")
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
			if selectedRow >= len(rows) {
				selectedRow = len(rows) - 1
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
