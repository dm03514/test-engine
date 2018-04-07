package main

import (
	"context"
	"flag"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine"
	"github.com/dm03514/test-engine/transcons"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func executeTest(fp string) error {
	ar, err := actions.NewRegistry()
	if err != nil {
		log.Panic(err)
	}

	tcr, err := transcons.NewRegistry()
	if err != nil {
		log.Panic(err)
	}

	dir, file := filepath.Split(fp)

	fl, err := engine.NewFileLoader(dir, ar, tcr, engine.NewDefaultFactory())
	if err != nil {
		log.Panic(err)
	}

	engine, err := fl.Load(file)

	if err != nil {
		log.Panic(err)
	}

	err = engine.Run(context.Background())
	return err
}

func main() {
	var fp = flag.String("test", "", "test to execute")
	flag.Parse()

	err := executeTest(*fp)

	if err != nil {
		log.Panic(err)
	}

	log.WithFields(log.Fields{
		"component": "test-executor.main",
	}).Info("SUCCESS")
}
