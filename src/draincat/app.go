package main

import (
	"bufio"
	"log"
	"net/http"
	"github.com/bmizerany/lpx"
)

func logsHandler(w http.ResponseWriter, r *http.Request) {
	lp := lpx.NewReader(bufio.NewReader(r.Body))
	for lp.Next() {
		hdr := lp.Header()
		data := lp.Bytes()
		err := Insert(
			hdr.PrivalVersion, hdr.Time, hdr.Hostname, hdr.Name, hdr.Procid, hdr.Msgid, data)
		if err != nil {
			// Fail abruptly as we do not know the appropriate response here.
			log.Fatalf("Database insert error: %v\n", err)
		}
	}
}

func main() {
	log.Printf("Running app on port %v\n", config.Port)
	http.HandleFunc("/logs", logsHandler)
	http.ListenAndServe("0.0.0.0:"+config.Port, nil)
}
