package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/lpx"
	"log"
	"net/http"
)

type LogLine struct {
	PrivalVersion, Time, HostName, Name, ProcID, MsgID, Data string
}

func NewLogLineFromLpx(lp *lpx.Reader) *LogLine {
		hdr := lp.Header()
		data := lp.Bytes()
		return &LogLine{
			string(hdr.PrivalVersion),
			string(hdr.Time),
			string(hdr.Hostname),
			string(hdr.Name),
			string(hdr.Procid),
			string(hdr.Msgid),
			string(data),
		}
}

var logsCh chan *LogLine

func receiveLogs() {
	for line := range logsCh {
		err := handleLog(line)
		if err != nil {
			log.Fatalf("Error handling log: %v", err)
		}
	}
}
func handleLog(line *LogLine) error {
	var err error
	if config.Json {
		data, err := json.Marshal(&line)
		if err != nil {
			log.Fatalf("JSON error: %v", err)
		}
		_, err = fmt.Println(string(data))
	} else {
		_, err = fmt.Printf("==> %v, %v, %v, %v, %v, %v, %v",
			line.PrivalVersion, line.Time, line.HostName, line.Name,
			line.ProcID, line.MsgID, line.Data)
	}
	return err
}

func routeLogs(w http.ResponseWriter, r *http.Request) {
	lp := lpx.NewReader(bufio.NewReader(r.Body))
	for lp.Next() {
		logsCh <- NewLogLineFromLpx(lp)
	}
}

func main() {
	logsCh = make(chan *LogLine)
	go receiveLogs()

	http.HandleFunc("/logs", routeLogs)
	http.ListenAndServe("0.0.0.0:"+config.Port, nil)
}
