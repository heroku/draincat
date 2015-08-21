package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/lpx"
	"log"
	"net/http"
	"os"
)

type LogLine struct {
	PrivalVersion string `json:"priv"`
	Time          string `json:"time"`
	HostName      string `json:"hostname"`
	Name          string `json:"name"`
	ProcID        string `json:"procid"`
	MsgID         string `json:"msgid"`
	Data          string `json:"data"`
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

var logsCh chan []*LogLine

const LOGSCH_BUFFER = 100

func receiveLogs() {
	for logs := range logsCh {
		for _, line := range logs {
			err := handleLog(line)
			if err != nil {
				log.Fatalf("Error handling log: %v", err)
			}
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
	logs := make([]*LogLine, 0)
	for lp.Next() {
		logs = append(logs, NewLogLineFromLpx(lp))
	}
	logsCh <- logs
}

func main() {
	if config.Port == 0 {
		fmt.Fprintln(os.Stderr, "err: invalid port")
		os.Exit(2)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)

	logsCh = make(chan []*LogLine, LOGSCH_BUFFER)
	go receiveLogs()

	http.HandleFunc("/logs", routeLogs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http server failure: %v\n", err)
		os.Exit(2)
	}
}
