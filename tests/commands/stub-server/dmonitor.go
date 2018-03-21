package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type DmCreated struct {
	ID string

	mu                   sync.Mutex
	numPolledCreated     int
	currPollCreatedCount int
}

func (d *DmCreated) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")
}

func (d *DmCreated) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(DmCreated{ID: "ID-CREATED"})
}

func (d *DmCreated) analysisComplete(w http.ResponseWriter, r *http.Request) {
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

func NewDmcreated(numPolledCreated int) *DmCreated {
	dm := &DmCreated{
		numPolledCreated: numPolledCreated,
	}

	http.HandleFunc("/dmonitor/delete", dm.delete)
	http.HandleFunc("/dmonitor/create", dm.create)
	http.HandleFunc("/dmonitor/analysis_complete", dm.analysisComplete)

	return dm
}
