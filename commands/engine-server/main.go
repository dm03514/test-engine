package main

import (
	"flag"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
)

func main() {
	var testsDir = flag.String("testDir", "", "Path to serve tests from")
	flag.Parse()

	log.Infof("testsDir: %q", *testsDir)

	ar, err := actions.NewActionRegistry()
	if err != nil {
		log.Panic(err)
	}

	tcr, err := transcons.NewTransConsRegistry()
	if err != nil {
		log.Panic(err)
	}

	loader, err := engine.NewFileLoader(*testsDir, ar, tcr)
	if err != nil {
		log.Panic(err)
	}

	s, err := engine.NewHTTPExecutor(
		engine.NewLoaders(
			loader,
		),
	)
	if err != nil {
		log.Panic(err)
	}

	s.ListenAndServer()
}
