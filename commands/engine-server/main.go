package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

func main() {
	var testsDir = flag.String("testDir", "", "Path to serve tests from")
	flag.Parse()

	log.Infof("testsDir: %q", *testsDir)

	// start the engine server
}
