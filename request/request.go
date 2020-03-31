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

type EditArgs struct {
    Level int
}

type Request struct {
	Cmd      Cmd
	NotesDir string
    EditArgs *EditArgs
}

func bindSharedArgs(fs *flag.FlagSet, r *Request) {
	fs.String(r.NotesDir, "path", "path to notes directory")
}

func bindCommandArgs(fs *flag.FlagSet, r *Request) {
    if r.Cmd == EDIT {
        r.EditArgs = &EditArgs{}
        fs.IntVar(&r.EditArgs.Level, "level", 0, "level")
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
		bindSharedArgs(flagSet, r)
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
