package engine

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

type Loaders struct {
	ls []Loader
}

// Load by iterating through all loaders, returning
// the first one that matches or error
func (l Loaders) Load(name string) (*Engine, error) {
	for _, loader := range l.ls {
		e, err := loader.Load(name)
		if err != nil {
			return nil, err
		}

		return e, nil
	}

	return nil, fmt.Errorf("No engine found matching %q", name)
}

func NewLoaders(ls ...Loader) Loaders {
	return Loaders{ls: ls}
}

type Loader interface {
	Load(name string) (*Engine, error)
}

type FileLoader struct {
	Dir               string
	actionRegistry    ActionRegistry
	transConsRegistry TransConsRegistry
	engineFactory     Factory
}

// Load the test from the Dir matching the name
func (fl FileLoader) Load(name string) (*Engine, error) {
	p := filepath.Join(fl.Dir, name)
	log.Infof("engine.Load(%q) from %q", name, p)
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	engine, err := NewFromYaml(
		content,
		fl.actionRegistry,
		fl.transConsRegistry,
		fl.engineFactory,
	)
	return engine, err
}

func NewFileLoader(dir string, ar ActionRegistry, tcr TransConsRegistry, ef Factory) (FileLoader, error) {
	return FileLoader{
		Dir:               dir,
		actionRegistry:    ar,
		transConsRegistry: tcr,
		engineFactory:     ef,
	}, nil
}

type MemoryLoader struct {
	m map[string]*Engine
}

func (ml *MemoryLoader) Load(name string) (*Engine, error) {
	e, ok := ml.m[name]
	if !ok {
		return nil, fmt.Errorf("Engine %q not found in %+v", name, ml.m)
	}
	return e, nil
}
