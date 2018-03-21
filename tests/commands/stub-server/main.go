package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	var port = flag.String("port", ":9999", "Stub Server Port")
	flag.Parse()

	NewDmcreated(5)
	NewGstreamer(5)

	s := &http.Server{
		Addr:         *port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
