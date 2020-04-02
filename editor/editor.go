package editor

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	DefaultEditor = "vim"
	MaxDuplicates = 10000
)

func getPath(dir string, name string, duplicates int) string {
	if duplicates == 0 {
		return fmt.Sprintf("%s/%s.md", dir, name)
	}
	return fmt.Sprintf("%s/%s(%d).md", dir, name, duplicates+1)
}

func Edit(path string) error {
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

func CreateAndEdit(dir string, name string) error {
	duplicates := 0

	var path string
	var err error
	for {
		path = getPath(dir, name, duplicates)
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

	return Edit(path)
}
