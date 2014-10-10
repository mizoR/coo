package tee

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"time"
)

type Command struct {
	Stdin        *os.File
	Stdout       *os.File
	Files        [](*os.File)
	OptTimestamp bool
}

func NewCommand() *Command {
	cmd := &Command{
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		OptTimestamp: false,
		Files:        [](*os.File){},
	}
	return cmd
}

func (cmd Command) writeFiles(buf []byte) {
	for _, file := range cmd.Files {
		file.Write(buf)
	}
}

func (cmd Command) bumpTime() {
	timestamp := time.Now().Format("[2006-01-02 15:04:05 +MST] ")
	cmd.writeFiles(([]byte)(timestamp))
}

func (cmd Command) write(buf []byte) {
	cmd.Stdout.Write(buf)
	cmd.writeFiles(buf)
}

func (cmd Command) Run() error {
	var err error

	bufsize := 8 * 1024
	buf := make([]byte, bufsize, bufsize)

	reader := bufio.NewReader(cmd.Stdin)

	if cmd.OptTimestamp {
		linehead := true
		for {
			var n int
			if n, err = reader.Read(buf); err != nil {
				break
			}
			b := bytes.NewBuffer(buf[0:n])

			var line []byte
			for {
				if line, err = b.ReadBytes((byte)('\n')); err != nil {
					break
				}
				if linehead {
					cmd.bumpTime()
				} else {
					linehead = true
				}
				cmd.write(line)
			}

			if err == io.EOF {
				if linehead {
					cmd.bumpTime()
					linehead = false
				}
				cmd.write(line)
			}
		}
	} else {
		for {
			var n int
			if n, err = reader.Read(buf); err != nil {
				break
			}
			cmd.write(buf[0:n])
		}
	}

	if err == io.EOF {
		return nil
	}

	return err
}
