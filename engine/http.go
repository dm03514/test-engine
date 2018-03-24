package engine

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html"
	"net/http"
	"time"
)

// HTTPExecutor can load tests and execute them through REST interface!
type HTTPExecutor struct {
	Addr string

	Loaders Loaders
}

// Execute a test
func (he HTTPExecutor) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	testName := r.FormValue("test")

	e, err := he.Loaders.Load(testName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	go func() {
		err := e.Run(context.Background())
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("SUCCESS!")

	}()

	fmt.Fprintf(w, "Test: %q started", html.EscapeString(testName))
}

// ListenAndServe tests through http REST interface
func (he HTTPExecutor) ListenAndServe() {
	s := &http.Server{
		Addr:         he.Addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

// RegisterHandlers registers all routes
func (he HTTPExecutor) RegisterHandlers() {
	http.HandleFunc("/execute", he.Execute)
}

// NewHTTPExecutor initializes and creates http executor server
func NewHTTPExecutor(loaders Loaders) (HTTPExecutor, error) {
	return HTTPExecutor{
		Addr:    ":8080",
		Loaders: loaders,
	}, nil
}
