package manager

import (
	"bufio"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	maxNonMatchPenalty = -100.0
	maxMatchScore      = 250.0
	dateWeight         = 10.0
	datePenaltyCeil    = 100.0
	secondsInDay       = 86400.0
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

	datePenalty := (curTime.Sub(note.ModTime).Seconds() / secondsInDay) * dateWeight
	if datePenalty > datePenaltyCeil {
		datePenalty = datePenaltyCeil
	}

	return matchingScore - datePenalty
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

func SortNotes(notes []Note, searchKey string) {
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

func (m *Manager) getTags(f os.FileInfo) ([]string, error) {
	tags := []string{}

	file, err := os.Open(m.getPath(f.Name()))
	if err != nil {
		return tags, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := strings.TrimSpace(scanner.Text())
	err = scanner.Err()
	if err != nil {
		return tags, err
	}

	if len(firstLine) < 2 || firstLine[0] != '[' || firstLine[len(firstLine)-1] != ']' {
		return tags, nil
	}
	listItems := strings.Split(firstLine[1:len(firstLine)-1], ",")
	for _, item := range listItems {
		item = strings.TrimSpace(item)
		if len(item) > 0 && item[0] == '#' {
			tags = append(tags, item[1:])
		}
	}

	return tags, nil
}

func arraysOverlap(a []string, b []string, caseSensitive bool) bool {
	for _, sa := range a {
		for _, sb := range b {
			if strings.ToLower(sa) == strings.ToLower(sb) {
				return true
			}
		}
	}
	return false
}

func (m *Manager) ListNotes(tags []string) ([]Note, error) {
	notes := []Note{}

	files, err := ioutil.ReadDir(m.Dir)
	if err != nil {
		return notes, err
	}

	for _, f := range files {
		n := f.Name()
		if len(n) >= 3 && n[len(n)-3:] == ".md" {
			fileTags, err := m.getTags(f)
			if err != nil {
				return notes, err
			}

			// TODO: We are doing an OR here, we should also support AND
			if arraysOverlap(tags, fileTags, false) {
				notes = append(notes, Note{
					n[:len(n)-3],
					fileTags,
					m.getPath(f.Name()),
					f.ModTime(),
				})
			}
		}
	}

	return notes, nil
}