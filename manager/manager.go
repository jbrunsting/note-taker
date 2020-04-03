package manager

import "time"

type Note struct {
	Title   string
	ModTime time.Time
}

type Manager struct {
	Dir string
}
