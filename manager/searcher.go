package manager

import (
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

type Note struct {
	Title   string
	ModTime time.Time
}

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
