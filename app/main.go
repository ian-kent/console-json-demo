package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type logMessage struct {
	Date  time.Time              `json:"date"`
	Level string                 `json:"level"`
	Event string                 `json:"event"`
	Args  map[string]interface{} `json:"args"`
}

var levels = []string{"INFO", "DEBUG", "WARN", "TRACE", "ERROR"}
var events = []string{"Cake", "Pizza", "Coffee"}
var argKeys = []string{"cows", "ducks", "sheep"}

func main() {
	var buf *bytes.Buffer
	for {
		// wait some random time
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

		// create some random args
		args := make(map[string]interface{})
		for i := 0; i < rand.Intn(3); i++ {
			args[argKeys[rand.Intn(len(argKeys))]] = rand.Intn(500)
		}

		// write a log message
		b, _ := json.Marshal(&logMessage{
			Date:  time.Now(),
			Level: levels[rand.Intn(len(levels))],
			Event: events[rand.Intn(len(events))],
			Args:  args,
		})
		buf = new(bytes.Buffer)
		json.Compact(buf, b)
		buf.WriteByte('\n')
		os.Stdout.Write(buf.Bytes())
	}
}
