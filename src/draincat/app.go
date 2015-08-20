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

func handleLog(line LogLine) error {
	if config.Json {
		data, err := json.Marshal(&line)
		if err != nil {
			log.Fatalf("JSON error: %v", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Printf("==> %v, %v, %v, %v, %v, %v, %v",
			line.PrivalVersion, line.Time, line.HostName, line.Name,
			line.ProcID, line.MsgID, line.Data)
	}
	return nil
}

func routeLogs(w http.ResponseWriter, r *http.Request) {
	lp := lpx.NewReader(bufio.NewReader(r.Body))
	for lp.Next() {
		hdr := lp.Header()
		data := lp.Bytes()
		err := handleLog(LogLine{
			string(hdr.PrivalVersion),
			string(hdr.Time),
			string(hdr.Hostname),
			string(hdr.Name),
			string(hdr.Procid),
			string(hdr.Msgid),
			string(data),
		})
		if err != nil {
			// Fail abruptly as we do not know the appropriate response here.
			log.Fatalf("Failed to handle a log line: %v\n", err)
		}
	}
}

func main() {
	log.Printf("Running app on port %v\n", config.Port)
	http.HandleFunc("/logs", routeLogs)
	http.ListenAndServe("0.0.0.0:"+config.Port, nil)
}
