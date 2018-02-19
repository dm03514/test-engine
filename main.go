package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine"
	"github.com/dm03514/test-engine/transcons"
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
	ar, err := actions.NewActionRegistry()
	if err != nil {
		panic(err)
	}
	tcr, err := transcons.NewTransConsRegistry()
	if err != nil {
		panic(err)
	}
	engine, err := engine.New(content, ar, tcr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", engine)

	err = engine.Run(context.Background())
	if err != nil {
		panic(err)
	}

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
	fmt.Println("SUCCESS")
}
