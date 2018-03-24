package engine

import (
	"fmt"
	"github.com/dm03514/test-engine/actions"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

// Loaders collection of loaders
type Loaders struct {
	ls []loader
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

// NewLoaders creates a new collection of loaders
func NewLoaders(ls ...loader) Loaders {
	return Loaders{ls: ls}
}

type loader interface {
	Load(name string) (*Engine, error)
}

// FileLoader contains information on how /where to load files from
type FileLoader struct {
	Dir               string
	actionRegistry    actions.Registry
	transConsRegistry transConsRegistry
	engineFactory     Factory
}

// Load the test from the Dir matching the name
func (fl FileLoader) Load(name string) (*Engine, error) {
	p := filepath.Join(fl.Dir, name)

	log.WithFields(log.Fields{
		"component": "Fileloader.Load()",
		"filename":  name,
		"path":      p,
	}).Info("loading_test")

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

// NewFileLoader creates a file loader
func NewFileLoader(dir string, ar actions.Registry, tcr transConsRegistry, ef Factory) (FileLoader, error) {
	return FileLoader{
		Dir:               dir,
		actionRegistry:    ar,
		transConsRegistry: tcr,
		engineFactory:     ef,
	}, nil
}

// MemoryLoader stores engines by identifier
type MemoryLoader struct {
	m map[string]*Engine
}

// Load an engine based on its name
func (ml *MemoryLoader) Load(name string) (*Engine, error) {
	e, ok := ml.m[name]
	if !ok {
		return nil, fmt.Errorf("Engine %q not found in %+v", name, ml.m)
	}
	return e, nil
}
