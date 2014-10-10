package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mizoR/coma/ssh"
)

func showHelp() {
	os.Stderr.WriteString("Usage: coma host [command]\n")
}

func main() {
	var st int = 1
	var err error
	var args []string

	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(st)
	}()

	parser := flags.NewParser(nil, flags.PrintErrors)
	args, err = parser.Parse()

	if err != nil || len(args) == 0 {
		showHelp()
		return
	}

	cmd := ssh.NewCommand()
	cmd.Host = args[0]
	cmd.Command = args[1:]
	cmd.Run()

	if err = cmd.Run(); err != nil {
		return
	}

	st = 0
}
