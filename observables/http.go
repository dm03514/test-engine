package observables

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

// configure to listen on a port
// Push all requests to channel
// CLEANUP ON TEST FINISH THIS NEEDS TO CLOSE AND CLEANUP CHANNEL

// HTTPRequestEvent is emitted whenever a request is received
type HTTPRequestEvent struct {
	*http.Request
}

// HTTP Listens and emits all HTTP Requests on a specific Addr
type HTTP struct {
	Addr string
}

// RunUntilContextDone listens for HTTP traffic and emits all requests,
// until the context is Done
func (h HTTP) RunUntilContextDone(ctx context.Context) <-chan ObservableEvent {
	reqs := make(<-chan ObservableEvent, 1000)
	/*
		s := &http.Server{
			Addr:         h.Addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		}
	*/

	go func() {

		// listen for context done
		// start the server
		// when context done is complete force a shutdown and close the channel

		/*
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				// Error starting or closing listener:
				log.Printf("HTTP server ListenAndServe: %v", err)
			}
		*/

	}()

	return reqs
}

// NewHTTPObservableFromMap initializes an HTTP
func NewHTTPObservableFromMap(m map[string]interface{}) (Observable, error) {
	var h HTTP
	err := mapstructure.Decode(m, &h)
	return h, err
}
