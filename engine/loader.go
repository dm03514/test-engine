package engine

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

type Loader interface {
	Load(name string) *Engine
}

type FileLoader struct {
	Dir string
	ar  ActionRegistry
	tcr TransConsRegistry
}

// Load the test from the Dir matching the name
func (fl FileLoader) Load(name string) (*Engine, error) {
	p := filepath.Join(fl.Dir, name)
	log.Infof("eninge.Load(%q) from %q", name, p)
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	engine, err := New(content, fl.ar, fl.tcr)
	return engine, err
}

func NewFileLoader(dir string, ar ActionRegistry, tcr TransConsRegistry) (FileLoader, error) {
	return FileLoader{
		Dir: dir,
		ar:  ar,
		tcr: tcr,
	}, nil
}
