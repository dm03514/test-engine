package engine

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html"
	"net/http"
	"time"
)

type HTTPOpt string

type httpExecutor struct {
	Addr string

	Loaders Loaders
}

func (he httpExecutor) execute(w http.ResponseWriter, r *http.Request) {
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

func (he httpExecutor) ListenAndServe() {
	http.HandleFunc("/execute", he.execute)

	s := &http.Server{
		Addr: he.Addr,
		// Handler:        myHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

func NewHTTPExecutor(loaders Loaders, opts ...HTTPOpt) (httpExecutor, error) {
	return httpExecutor{
		Addr:    ":8080",
		Loaders: loaders,
	}, nil
}
