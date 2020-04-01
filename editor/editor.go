package editor

import (
	"os"
	"os/exec"
)

const (
	DefaultEditor = "vim"
)

func CreateAndEdit(filename string, header string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return err
	}
	if fi.Size() == 0 {
		file.Write([]byte(header))
	}

	file.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
