package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine"
	"github.com/dm03514/test-engine/fulfillment"
	"github.com/dm03514/test-engine/transcons"
	"time"
)

func main() {
	var fp = flag.String("test", "", "test to execute")
	flag.Parse()
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

	t := engine.Test{
		States: []engine.State{
			fulfillment.NoopFulillment{
				Action: actions.Subprocess{
					"echo",
					[]string{"hello world!"},
				},
				Conditions: transcons.Conditions{
					[]transcons.TransCon{
						transcons.IntEqual{
							UsingProperty: "returncode",
							ToEqual:       0,
						},
					},
				},
			},
		},
		Timeout: time.Duration(1 * time.Minute),
	}

	e := engine.Engine{
		Test: t,
	}
	ctx := context.Background()
	err := e.Run(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("SUCESS!")

}
