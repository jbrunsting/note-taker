package manager

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	DefaultEditor = "vim"
	MaxDuplicates = 10000
)

func (m *Manager) getPath(name string) string {
	return fmt.Sprintf("%s/%s", m.Dir, name)
}

func (m *Manager) getFileName(name string, extension string, duplicates int) string {
	if duplicates == 0 {
		return fmt.Sprintf("%s.%s", name, extension)
	}
	return fmt.Sprintf("%s(%d).%s", name, duplicates+1, extension)
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
	path := m.getPath(m.getFileName(name, "md", 0))
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return edit(m.getPath(m.getFileName(name, "md", 0)))
}

func (m *Manager) Move(src string, name string, extension string) error {
	duplicates := 0

	var path string
	var err error
	for {
		path = m.getPath(m.getFileName(name, extension, duplicates))
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

	return os.Rename(src, path)
}

func (m *Manager) CreateAndEdit(name string, header string) error {
	duplicates := 0

	var path string
	var err error
	for {
		path = m.getPath(m.getFileName(name, "md", duplicates))
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

func (m *Manager) ReadNote(note *Note) ([]string, error) {
	content, err := ioutil.ReadFile(note.Path)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(content), "\n"), nil
}

func (m *Manager) Delete(name string) error {
	return os.Remove(m.getPath(m.getFileName(name, "md", 0)))
}

func (m *Manager) ViewAll(notes []Note) error {
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

	for i, note := range notes {
		noteFile, err := os.Open(note.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer noteFile.Close()

		pre := fmt.Sprintf("# %s\n\n", note.Title)
		if i != 0 {
			pre = "\n" + pre
		}
		_, err = file.WriteString(pre)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(noteFile)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			if len(line) > 0 && line[0] == '#' {
				// Indent markdown headers by 1
				line = "#" + line
			}

			_, err = file.WriteString(line)
			if err != nil {
				return err
			}
		}
	}

	return edit(path)
}
