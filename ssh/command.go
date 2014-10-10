package ssh

import (
	"os"
	"os/exec"
	"time"
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
	tee := exec.Command("coo-tee", cmd.Logfile, "-t", "-a")

	ssh.Stdin = os.Stdin
	ssh.Stderr = os.Stderr
	tee.Stdout = os.Stdout

	if tee.Stdin, err = ssh.StdoutPipe(); err != nil {
		return err
	}

	if err = ssh.Start(); err != nil {
		return err
	}

	if err = tee.Start(); err != nil {
		return err
	}

	if err = ssh.Wait(); err != nil {
		return err
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
