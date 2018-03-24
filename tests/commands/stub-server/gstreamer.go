package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// GStreamer contains metadata about streaming
type GStreamer struct {
	AccountID string

	mu                   sync.Mutex
	numPolledCreated     int
	currPollCreatedCount int
}

func (gs *GStreamer) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(GStreamer{AccountID: "ID-CREATED"})
}

func (gs *GStreamer) analysisComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	time.Sleep(1 * time.Second)
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if gs.currPollCreatedCount != gs.numPolledCreated {
		gs.currPollCreatedCount++
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "{\"results\":[]}")
	} else {
		gs.currPollCreatedCount = 0
		fmt.Fprintf(w, "{\"results\":[\"success\"]}")
	}
}

// NewGstreamer creates initializes a new streamer test server
func NewGstreamer(numPolledCreated int) *GStreamer {
	gs := &GStreamer{
		numPolledCreated: numPolledCreated,
	}

	http.HandleFunc("/gstreamer/account/create", gs.create)
	http.HandleFunc("/gstreamer/analysis_complete", gs.analysisComplete)

	return gs
}
