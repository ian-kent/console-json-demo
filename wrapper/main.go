package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mgutz/ansi"
)

type logMessage struct {
	Date  time.Time              `json:"date"`
	Level string                 `json:"level"`
	Event string                 `json:"event"`
	Args  map[string]interface{} `json:"args"`
}

type logWriter struct {
	buf *bytes.Buffer
}

func (w *logWriter) Write(b []byte) (n int, err error) {
	if w.buf == nil {
		w.buf = new(bytes.Buffer)
	}
	for _, c := range b {
		n++
		if c == '\n' {
			var msg logMessage
			err = json.Unmarshal(w.buf.Bytes(), &msg)
			if err != nil {
				panic(err)
			}
			w.buf = new(bytes.Buffer)
			var col = ansi.DefaultFG
			switch msg.Level {
			case "ERROR":
				col = ansi.LightRed
			case "TRACE":
				col = ansi.Blue
			case "INFO":
				col = ansi.Cyan
			case "WARN":
				col = ansi.Yellow
			case "DEBUG":
				col = ansi.Green
			}
			fmt.Fprintf(os.Stdout, "%s%s [%s] %s%s\n", col, msg.Date, msg.Level, msg.Event, ansi.DefaultFG)
			for k, v := range msg.Args {
				fmt.Fprintf(os.Stdout, "  => %s%s:%s %+v\n", ansi.LightBlack, k, ansi.DefaultFG, v)
			}
			continue
		}
		err = w.buf.WriteByte(c)
	}
	return
}

func main() {
	// create a command using passed in command line args
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	// redirect console stderr/stdin to child process
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	writer := &logWriter{}
	cmd.Stdout = writer

	// start the process
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// wait until it finishes
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
