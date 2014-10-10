package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"
	"github.com/mizoR/coo/tee"
	"github.com/mizoR/coo/usage"
)

type cmdOptions struct {
	OptHelp      bool `short:"h" long:"help" description:"Show this help message and exit"`
	OptAppend    bool `short:"a" long:"append" description:"Append the output to the files rather than overwriting them."`
	OptTimestamp bool `short:"t" long:"timestamp" description:"Record the time to the output."`
}

func showHelp() {
	os.Stderr.WriteString("Usage: coo-tee [options] [FILE]\n")

	t := reflect.TypeOf(cmdOptions{})
	usage.Show(t)
}

func main() {
	var err error
	var st int = 1

	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(st)
	}()

	opts := &cmdOptions{}
	p := flags.NewParser(opts, flags.PrintErrors)
	files, err := p.Parse()

	if err != nil {
		showHelp()
		return
	}

	if opts.OptHelp {
		st = 0
		showHelp()
		return
	}

	cmd := tee.NewCommand()
	cmd.OptTimestamp = opts.OptTimestamp

	perm := (os.O_WRONLY | os.O_CREATE)
	if opts.OptAppend {
		perm = (perm | os.O_APPEND)
	}

	for _, file := range files {
		var f *os.File
		if f, err = os.OpenFile(file, perm, 0644); err != nil {
			return
		}
		defer f.Close()
		cmd.Files = append(cmd.Files, f)
	}

	if err = cmd.Run(); err != nil {
		return
	}

	st = 0
}
