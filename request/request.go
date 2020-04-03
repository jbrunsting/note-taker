package request

import (
	"flag"
	"log"
	"os"
)

type Cmd int

const (
	NEW Cmd = iota
	EDIT
)

const (
	DefaultNoteTitle = "Untitled"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type NewArgs struct {
	Title string
	Tags  arrayFlags
}

type EditArgs struct {
	Title string
	Tags  arrayFlags
}

type Request struct {
	Cmd      Cmd
	NotesDir string
	NewArgs  *NewArgs
	EditArgs *EditArgs
}

func bindSharedArgs(fs *flag.FlagSet, r *Request) {
	fs.StringVar(&r.NotesDir, "path", "", "path to notes directory")
}

func bindCommandArgs(fs *flag.FlagSet, r *Request) {
	if r.Cmd == NEW {
		r.NewArgs = &NewArgs{}
		fs.StringVar(&r.NewArgs.Title, "title", DefaultNoteTitle, "the title of the note")
		fs.Var(&r.NewArgs.Tags, "tags", "the tags for the note")
	} else if r.Cmd == EDIT {
		r.EditArgs = &EditArgs{}
		fs.StringVar(&r.EditArgs.Title, "title", "", "the title of the note")
		fs.Var(&r.EditArgs.Tags, "tags", "the tags for the note")
	}
}

func CreateRequest() Request {
	if len(os.Args) < 2 {
		log.Fatalf("TODO: Print subcommands\n")
		os.Exit(1)
	}

	cmds := make(map[string]Cmd)
	cmds["new"] = NEW
	cmds["edit"] = EDIT

	flagSets := make(map[Cmd]*flag.FlagSet)
	flagSets[NEW] = flag.NewFlagSet("new", flag.ExitOnError)
	flagSets[EDIT] = flag.NewFlagSet("edit", flag.ExitOnError)

	var r Request
	for _, flagSet := range flagSets {
		bindSharedArgs(flagSet, &r)
	}

	if cmd, ok := cmds[os.Args[1]]; ok {
		r.Cmd = cmd
		bindCommandArgs(flagSets[cmd], &r)
		flagSets[cmd].Parse(os.Args[2:])
	} else {
		log.Printf("TODO: Print proper error\n")
		os.Exit(1)
	}

	return r
}
