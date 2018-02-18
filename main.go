package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	var fp = flag.String("test", "", "test to execute")
	flag.Parse()
	content, err := ioutil.ReadFile(*fp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", content)
	/*
		result, err = engine.New(
			parser.NewYAML(
				parser.NewEnvVar(
					parser.NewUnqique(
						LoadFile(path)
					)
				)
			) -> TestStateMachine
		).Run()
	*/

	// Test State Machine
	fmt.Println(*fp)
}
