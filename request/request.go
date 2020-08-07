package request

import (
	"flag"
	"log"
	"os"
)

type Cmd int

const (
	NEW Cmd = iota
	MV
	EDIT
	DELETE
	CONCAT
	FIND
	HTML
	GIT
	PUSH
	INIT_REPO
)

const (
	DefaultNoteTitle = "Untitled"
)

type NewArgs struct {
	Title string
	Tags  ArrayFlags
}

type MvArgs struct {
	Title string
	Src   string
}

type EditArgs struct {
	Title string
	Tags  ArrayFlags
}

type DeleteArgs struct {
	Title string
}

type ConcatArgs struct {
	Tags ArrayFlags
}

type FindArgs struct {
	Tags ArrayFlags
}

type HtmlArgs struct {
	Tags ArrayFlags
	File string
}

type Request struct {
	Cmd        Cmd
	Args       []string
	NotesDir   string
	NewArgs    *NewArgs
	MvArgs     *MvArgs
	EditArgs   *EditArgs
	DeleteArgs *DeleteArgs
	ConcatArgs *ConcatArgs
	FindArgs   *FindArgs
	HtmlArgs   *HtmlArgs
}

func bindSharedArgs(fs *flag.FlagSet, r *Request) {
	fs.StringVar(&r.NotesDir, "path", "", "path to notes directory")
}

func bindCommandArgs(fs *flag.FlagSet, r *Request) {
	if r.Cmd == NEW {
		r.NewArgs = &NewArgs{}
		fs.StringVar(&r.NewArgs.Title, "title", DefaultNoteTitle, "the title of the note")
		fs.Var(&r.NewArgs.Tags, "tags", "the tags for the note")
	} else if r.Cmd == MV {
		r.MvArgs = &MvArgs{}
		fs.StringVar(&r.MvArgs.Title, "title", DefaultNoteTitle, "the new title of the file")
		fs.StringVar(&r.MvArgs.Src, "src", "", "the path to the file")
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
	} else if r.Cmd == FIND {
		r.FindArgs = &FindArgs{}
		fs.Var(&r.FindArgs.Tags, "tags", "the tags for the note")
	} else if r.Cmd == HTML {
		r.HtmlArgs = &HtmlArgs{}
		fs.Var(&r.HtmlArgs.Tags, "tags", "the tags for the note")
		fs.StringVar(&r.HtmlArgs.File, "file", "", "the file to store the html output")
	}
}

func RequestFromArgs() Request {
	cmds := make(map[string]Cmd)
	cmds["new"] = NEW
	cmds["mv"] = MV
	cmds["edit"] = EDIT
	cmds["delete"] = DELETE
	cmds["concat"] = CONCAT
	cmds["find"] = FIND
	cmds["html"] = HTML
	cmds["git"] = GIT
	cmds["push"] = PUSH
	cmds["init-repo"] = INIT_REPO

	keys := []string{}
	for k := range cmds {
		keys = append(keys, k)
	}
	if len(os.Args) < 2 {
		log.Fatalf("Must provide one of: %v\n", keys)
		os.Exit(1)
	}

	flagSets := make(map[Cmd]*flag.FlagSet)
	for s, e := range cmds {
		flagSets[e] = flag.NewFlagSet(s, flag.ExitOnError)
	}

	var r Request
	for _, flagSet := range flagSets {
		bindSharedArgs(flagSet, &r)
	}

	if cmd, ok := cmds[os.Args[1]]; ok {
		r.Args = os.Args[2:]
		r.Cmd = cmd
		bindCommandArgs(flagSets[cmd], &r)
		flagSets[cmd].Parse(os.Args[2:])
	} else {
		log.Fatalf("Unknown command '%s', must provide one of: %v\n", os.Args[1], keys)
		os.Exit(1)
	}

	return r
}
