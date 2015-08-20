package main

import (
	"bufio"
	"log"
	"net/http"
	"github.com/bmizerany/lpx"
	"fmt"
)

func handleLog(privalVersion, time, hostname, name, procid, msgid, data string) error {
	fmt.Printf("==> %v, %v, %v, %v, %v, %v, %v",
		privalVersion, time, hostname, name,
		procid, msgid, data)
	return nil
}

func routeLogs(w http.ResponseWriter, r *http.Request) {
	lp := lpx.NewReader(bufio.NewReader(r.Body))
	for lp.Next() {
		hdr := lp.Header()
		data := lp.Bytes()
		err := handleLog(
			string(hdr.PrivalVersion),
			string(hdr.Time),
			string(hdr.Hostname),
			string(hdr.Name),
			string(hdr.Procid),
			string(hdr.Msgid),
			string(data))
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
