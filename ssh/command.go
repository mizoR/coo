package ssh

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/mizoR/coo/tee"
)

type Command struct {
	Stdin   *os.File
	Stdout  *os.File
	Host    string
	Command [](string)
	Logfile string
}

func NewCommand() *Command {
	cmd := &Command{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Command: [](string){},
	}
	return cmd
}

func (cmd Command) createLogdir() (string, error) {
	var err error

	sshdir := os.Getenv("HOME") + "/.ssh"
	logdir := sshdir + "/transcripts/" + cmd.Host

	if _, err = os.Stat(sshdir); os.IsNotExist(err) {
		return logdir, err
	}

	if err = os.MkdirAll(logdir, 0700); err != nil {
		return logdir, err
	}

	return logdir, nil
}

func (cmd Command) execute() error {
	var err error
	args := append([]string{cmd.Host}, cmd.Command...)
	ssh := exec.Command("/usr/bin/ssh", args...)
	tee := tee.NewTee(cmd.Logfile, true, true)

	ssh.Stdin = os.Stdin
	ssh.Stderr = os.Stderr

	var reader io.Reader
	reader, err = ssh.StdoutPipe()
	if err != nil {
		return err
	}

	sshCh := make(chan error)
	go func() {
		sshCh <- ssh.Run()
	}()

	teeCh := tee.WriteBackground(reader)

	return wait((<-chan error)(sshCh), teeCh)

}

func wait(sshCh, teeCh <-chan error) error {
	var sshDone, teeDone bool
	for {
		if sshDone && teeDone {
			break
		}
		select {
		case err := <-sshCh:
			if err != nil {
				return err
			}
			sshDone = true
		case err := <-teeCh:
			if err != nil {
				return err
			}
			teeDone = true
		}

	}
	return nil
}

func (cmd Command) Run() error {
	var err error
	var logdir string

	if logdir, err = cmd.createLogdir(); err != nil {
		return err
	}

	cmd.Logfile = logdir + "/" + time.Now().Format("2006-01-02.txt")
	cmd.execute()

	return err
}
