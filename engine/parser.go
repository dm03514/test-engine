package engine

import (
	"fmt"
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/fulfillment"
	"github.com/dm03514/test-engine/transcons"
	"github.com/go-yaml/yaml"
	"github.com/mitchellh/mapstructure"
	"time"
)

type intermediaryState struct {
	Name                 string
	Fulfillment          map[string]interface{}
	Action               map[string]interface{}
	TransitionConditions []map[string]interface{} `yaml:"transition_conditions"`
}

func (is intermediaryState) ParsedAction() (actions.Action, error) {
	var a actions.Action
	var err error
	actionType := is.Action["type"].(string)

	switch actionType {
	case "shell.Subprocess":
		var subprocess actions.Subprocess
		err = mapstructure.Decode(is.Action, &subprocess)
		a = subprocess
	default:
		err = fmt.Errorf("Unable to parse action type %s", actionType)
	}

	return a, err
}

func (is intermediaryState) ParsedTransCons() (transcons.Conditions, error) {
	var err error
	var parsedCondition transcons.TransCon
	conditions := []transcons.TransCon{}

	for _, tc := range is.TransitionConditions {
		tcType := tc["type"].(string)
		switch tcType {
		case "assertions.IntEqual":
			var intEqual transcons.IntEqual
			err = mapstructure.Decode(tc, &intEqual)
			parsedCondition = intEqual
		default:
			err = fmt.Errorf("Unable to parse condition %s", tcType)
		}

		if err != nil {
			return transcons.Conditions{}, err
		}

		conditions = append(conditions, parsedCondition)
	}

	return transcons.Conditions{
		Tcs: conditions,
	}, nil
}

func (is intermediaryState) State() (State, error) {
	// get the action

	// get all transition conditions

	// get the eventfulfillment
	action, err := is.ParsedAction()
	if err != nil {
		return nil, err
	}
	conditions, err := is.ParsedTransCons()
	if err != nil {
		return nil, err
	}
	return fulfillment.NoopFulillment{
		Action:     action,
		Conditions: conditions,
	}, nil
}

type intermediaryTest struct {
	Name    string
	Timeout int
	States  []intermediaryState
}

func (it intermediaryTest) TimeoutDuration() time.Duration {
	to := 60
	if it.Timeout != 0 {
		to = it.Timeout
	}
	return time.Duration(time.Duration(to) * time.Second)
}

func New(b []byte) (*Engine, error) {
	it := intermediaryTest{}
	err := yaml.Unmarshal(b, &it)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", it)
	states := []State{}
	for _, ps := range it.States {
		fmt.Printf("%+v\n", ps)
		s, err := ps.State()
		if err != nil {
			return nil, err
		}
		states = append(states, s)
	}

	return &Engine{
		Test: Test{
			Name:    it.Name,
			Timeout: it.TimeoutDuration(),
			States:  states,
		},
	}, nil
}
