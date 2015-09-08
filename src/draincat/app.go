package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/lpx"
	"github.com/docopt/docopt-go"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
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

var logsCh chan *LogLine

const LOGSCH_BUFFER = 100

func receiveLogs(useJson bool) {
	for line := range logsCh {
		err := handleLog(line, useJson)
		if err != nil {
			log.Fatalf("Error handling log: %v", err)
		}
	}
}

func handleLog(line *LogLine, useJson bool) error {
	var err error
	if useJson {
		data, err := json.Marshal(&line)
		if err != nil {
			log.Fatalf("JSON error: %v", err)
		}
		_, err = fmt.Println(string(data))
	} else {
		_, err = fmt.Printf("%v, %v, %v, %v, %v, %v, %v",
			line.PrivalVersion, line.Time, line.HostName, line.Name,
			line.ProcID, line.MsgID, line.Data)
	}
	return err
}

var latency *Latency

func routeLogs(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(os.Stderr, "DEBUG: in request\n")
	if latency != nil {
		ms := latency.Create()
		time.Sleep(time.Duration(ms) * time.Millisecond)
		fmt.Fprintf(os.Stderr, "DEBUG: introduced %v delay in this response\n", ms)
	} else {
		// fmt.Fprintf(os.Stderr, "DEBUG: no delay\n")
	}

	lp := lpx.NewReader(bufio.NewReader(r.Body))
	for lp.Next() {
		logsCh <- NewLogLineFromLpx(lp)
	}
}

func oops(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	usage := `draincat
Usage:
  draincat [-j] [-D PERC] -p PORT
Options:
  -p PORT --port=PORT                 HTTP port to listen
  -j --json                           Output log messages in JSON
  -D PERC --latency-percentiles=PERC  Handle responses with the given percentile delay
`

	arguments, err := docopt.Parse(usage, nil, true, "draincat", false)
	if err != nil {
		oops(err)
	}
	portString := arguments["--port"].(string)
	useJson := arguments["--json"].(bool)
	latencyPercentiles, ok := arguments["--latency-percentiles"].(string)
	if ok {
		latency, err = NewLatencyFromSpec(latencyPercentiles)
		if err != nil {
			oops(fmt.Errorf("Invalid latency spec (%v): %v", latencyPercentiles, err))
		} else {
			fmt.Fprintf(os.Stderr, "WARNING: running draincat with latency distribution: %+v\n", latency)
		}
	}

	port, err := strconv.Atoi(portString)
	if err != nil || port == 0 {
		oops(fmt.Errorf("err: invalid port %s\n", portString))
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)

	logsCh = make(chan *LogLine, LOGSCH_BUFFER)
	go receiveLogs(useJson)

	rand.Seed(time.Now().Unix())

	http.HandleFunc("/logs", routeLogs)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http server failure: %v\n", err)
		os.Exit(2)
	}
}
