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

type HttpExecutor struct {
	Addr string

	Loaders Loaders
}

func (he HttpExecutor) Execute(w http.ResponseWriter, r *http.Request) {
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

func (he HttpExecutor) ListenAndServe() {
	s := &http.Server{
		Addr:         he.Addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

func (he HttpExecutor) RegisterHandlers() {
	http.HandleFunc("/execute", he.Execute)
}

func NewHTTPExecutor(loaders Loaders, opts ...HTTPOpt) (HttpExecutor, error) {
	return HttpExecutor{
		Addr:    ":8080",
		Loaders: loaders,
	}, nil
}
