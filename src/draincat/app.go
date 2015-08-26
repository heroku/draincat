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

var randomDelay bool

func routeLogs(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(os.Stderr, "DEBUG: in request\n")
	if randomDelay {
		_ = randomSleep(250, 750)
		// fmt.Fprintf(os.Stderr, "DEBUG: introduced %v delay in this response\n", ms)
	} else {
		// fmt.Fprintf(os.Stderr, "DEBUG: no delay\n")
	}
	os.Stderr.Sync()

	lp := lpx.NewReader(bufio.NewReader(r.Body))
	for lp.Next() {
		logsCh <- NewLogLineFromLpx(lp)
	}
}

func randomSleep(start, rng int) time.Duration {
	ms := time.Duration(time.Duration(start)*time.Millisecond + time.Duration(rand.Intn(rng))*time.Millisecond)
	time.Sleep(ms)
	return ms
}

func main() {
	usage := `draincat
Usage:
  draincat [-j] [-D] -p PORT
Options:
  -p PORT --port=PORT    HTTP port to listen
  -j --json              Output log messages in JSON
  -D --random-delay      Handle responses with random delay
`

	arguments, err := docopt.Parse(usage, nil, true, "draincat", false)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	portString := arguments["--port"].(string)
	useJson := arguments["--json"].(bool)
	randomDelay = arguments["--random-delay"].(bool)
	fmt.Fprintf(os.Stderr, "DEBUG: Random delay? %v\n", randomDelay)

	port, err := strconv.Atoi(portString)
	if err != nil || port == 0 {
		fmt.Fprintf(os.Stderr, "err: invalid port %s\n", portString)
		os.Exit(2)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)

	logsCh = make(chan *LogLine, LOGSCH_BUFFER)
	go receiveLogs(useJson)

	http.HandleFunc("/logs", routeLogs)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http server failure: %v\n", err)
		os.Exit(2)
	}
}
