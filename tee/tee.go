package tee

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"time"
)

type tee struct {
	file          string
	appendToFile  bool
	withTimeStamp bool
}

func NewTee(file string, appendToFile, withTimeStamp bool) tee {
	return tee{
		file:          file,
		appendToFile:  appendToFile,
		withTimeStamp: withTimeStamp,
	}
}

func (t tee) WriteBackground(r io.Reader) <-chan error {
	ch := make(chan error)
	go func() {
		ch <- t.Write(r)
	}()
	return ch
}

func (t tee) Write(r io.Reader) error {

	// Open the log file.
	perm := (os.O_WRONLY | os.O_CREATE)
	if t.appendToFile {
		perm = (perm | os.O_APPEND)
	}
	f, err := os.OpenFile(t.file, perm, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Get MultiWriter that behaves UNIX tee.
	var log io.Writer
	if t.withTimeStamp {
		log = &timeStampWriter{w: f, linehead: true}
	} else {
		log = f
	}
	w := io.MultiWriter(os.Stdout, log)

	// Read line from the Reader.
	reader := bufio.NewReader(r)
	bufsize := 8 * 1024
	buf := make([]byte, bufsize, bufsize)

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
			w.Write(line)
		}

		if err == io.EOF {
			w.Write(line)
		}
	}

	return nil
}

type timeStampWriter struct {
	linehead bool
	w        io.Writer
}

func (w *timeStampWriter) Write(p []byte) (n int, err error) {
	var line []byte
	if w.linehead {
		timestamp := time.Now().Format("[2006-01-02 15:04:05 +MST] ")
		line = append([]byte(timestamp), p...)
	} else {
		line = p
	}
	w.linehead = bytes.HasSuffix(p, []byte("\n"))
	return w.w.Write(line)
}
