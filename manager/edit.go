package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	DefaultEditor = "vim"
	MaxDuplicates = 10000
)

func (m *Manager) getPath(name string) string {
	return fmt.Sprintf("%s/%s", m.Dir, name)
}

func (m *Manager) getFileName(name string, duplicates int) string {
	if duplicates == 0 {
		return fmt.Sprintf("%s.md", name)
	}
	return fmt.Sprintf("%s(%d).md", name, duplicates+1)
}

func edit(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *Manager) Edit(name string) error {
	return edit(m.getPath(m.getFileName(name, 0)))
}

func (m *Manager) CreateAndEdit(name string, header string) error {
	duplicates := 0

	var path string
	var err error
	for {
		path = m.getPath(m.getFileName(name, duplicates))
		_, err = os.Stat(path)
		if err != nil {
			break
		}

		duplicates += 1
		if duplicates > 10000 {
			return fmt.Errorf("All file names had conflicts")
		}
	}
	if !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(header)
	if err != nil {
		return err
	}

	return edit(path)
}

func (m *Manager) BulkEdit(names []Note) error {
	file, err := ioutil.TempFile(os.TempDir(), "*.md")
	if err != nil {
		return err
	}
	path := file.Name()
	defer os.Remove(path)

	err = file.Close()
	if err != nil {
		return err
	}

	file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("Test")
	if err != nil {
		return err
	}

	return edit(path)
}
