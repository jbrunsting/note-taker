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
	DELETE
	CONCAT
)

const (
	DefaultNoteTitle = "Untitled"
)

type NewArgs struct {
	Title string
	Tags  ArrayFlags
}

type EditArgs struct {
	Title string
	Tags  ArrayFlags
}

type DeleteArgs struct {
	Title string
}

type ConcatArgs struct {
	Title string
	Tags  ArrayFlags
}

type Request struct {
	Cmd        Cmd
	NotesDir   string
	NewArgs    *NewArgs
	EditArgs   *EditArgs
	DeleteArgs *DeleteArgs
	ConcatArgs *ConcatArgs
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
	} else if r.Cmd == CONCAT {
		r.ConcatArgs = &ConcatArgs{}
		fs.Var(&r.ConcatArgs.Tags, "tags", "the tags for the note")
	} else if r.Cmd == DELETE {
		r.DeleteArgs = &DeleteArgs{}
		fs.StringVar(&r.DeleteArgs.Title, "title", "", "the title of the note")
	}
}

func RequestFromArgs() Request {
	if len(os.Args) < 2 {
		log.Fatalf("TODO: Print subcommands\n")
		os.Exit(1)
	}

	cmds := make(map[string]Cmd)
	cmds["new"] = NEW
	cmds["edit"] = EDIT
	cmds["delete"] = DELETE
	cmds["concat"] = CONCAT

	flagSets := make(map[Cmd]*flag.FlagSet)
	flagSets[NEW] = flag.NewFlagSet("new", flag.ExitOnError)
	flagSets[EDIT] = flag.NewFlagSet("edit", flag.ExitOnError)
	flagSets[DELETE] = flag.NewFlagSet("delete", flag.ExitOnError)
	flagSets[CONCAT] = flag.NewFlagSet("concat", flag.ExitOnError)

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
