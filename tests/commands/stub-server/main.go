package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type dmCreated struct {
	ID string

	mu                   sync.Mutex
	numPolledCreated     int
	currPollCreatedCount int
}

func (d dmCreated) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
}

func (d dmCreated) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dmCreated{ID: "ID-CREATED"})
}

func (d *dmCreated) analysisComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	if d.currPollCreatedCount != d.numPolledCreated {
		d.currPollCreatedCount++
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "{\"results\":[]}")
		return
	} else {
		d.currPollCreatedCount = 0
		fmt.Fprintf(w, "{\"results\":[\"success\"]}")
		return
	}

}

func main() {
	var port = flag.String("port", ":9999", "Stub Server Port")
	flag.Parse()

	dm := &dmCreated{
		numPolledCreated: 5,
	}

	http.HandleFunc("/dmonitor/delete", dm.delete)
	http.HandleFunc("/dmonitor/create", dm.create)
	http.HandleFunc("/dmonitor/analysis_complete", dm.analysisComplete)

	s := &http.Server{
		Addr:         *port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
