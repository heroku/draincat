package main

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/heroku/drain"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func receiveLogs(d *drain.Drain, useJson bool) {
	for line := range d.Logs() {
		err := handleLog(line, useJson)
		if err != nil {
			log.Fatalf("Error handling log: %v", err)
		}
	}
}

func handleLog(line *drain.LogLine, useJson bool) error {
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

func latencyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if latency != nil {
			ms := latency.Create()
			time.Sleep(time.Duration(ms) * time.Millisecond)
			fmt.Fprintf(os.Stderr, "DEBUG: introduced %v delay in this response\n", ms)
		}
		h.ServeHTTP(w, r)
	})
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

	theDrain := drain.NewDrain()
	go receiveLogs(theDrain, useJson)

	rand.Seed(time.Now().Unix())

	logsHandler := latencyMiddleware(
		http.HandlerFunc(theDrain.LogsHandler))

	http.Handle("/logs", logsHandler)
	err = http.ListenAndServe(addr, nil)
	oops(err)
}
