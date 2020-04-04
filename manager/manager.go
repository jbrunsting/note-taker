package manager

import "time"

type Note struct {
	Id      int
	Title   string
	Tags    []string
	Path    string
	ModTime time.Time
}

type Manager struct {
	Dir string
}
