package engine

import (
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/fulfillment"
	"github.com/dm03514/test-engine/transcons"
	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
	"time"
)

type ActionRegistry interface {
	Load(am map[string]interface{}) (actions.Action, error)
}

type TransConsRegistry interface {
	Load(tcm map[string]interface{}) (transcons.TransCon, error)
}

type intermediaryState struct {
	Name                 string
	Fulfillment          map[string]interface{}
	Action               map[string]interface{}
	TransitionConditions []map[string]interface{} `yaml:"transition_conditions"`
}

func (is intermediaryState) ParsedTransCons(tcr TransConsRegistry) (transcons.Conditions, error) {
	var err error
	var parsedCondition transcons.TransCon
	conditions := make([]transcons.TransCon, len(is.TransitionConditions))

	for i, tc := range is.TransitionConditions {
		parsedCondition, err = tcr.Load(tc)

		if err != nil {
			return transcons.Conditions{}, err
		}
		conditions[i] = parsedCondition
	}

	return transcons.Conditions{
		Tcs: conditions,
	}, nil
}

func (is intermediaryState) State(ar ActionRegistry, tcr TransConsRegistry) (State, error) {
	action, err := ar.Load(is.Action)
	if err != nil {
		return nil, err
	}
	conditions, err := is.ParsedTransCons(tcr)
	if err != nil {
		return nil, err
	}

	fr, err := fulfillment.NewRegistry()
	if err != nil {
		return nil, err
	}
	return fr.Load(is.Fulfillment, is.Name, action, conditions)
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

func NewFromYaml(b []byte, ar ActionRegistry, tcr TransConsRegistry, f Factory) (*Engine, error) {
	it := intermediaryTest{}
	err := yaml.Unmarshal(b, &it)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"component": "NewFromYaml()",
		"test_name": it.Name,
	}).Debug("parsing_test")

	states := []State{}
	for _, ps := range it.States {

		log.WithFields(log.Fields{
			"component": "NewFromYaml()",
			"raw_state": ps,
		}).Debug("parsing_state")

		s, err := ps.State(ar, tcr)
		if err != nil {
			return nil, err
		}
		states = append(states, s)
	}

	return f.New(
		Test{
			Name:    it.Name,
			Timeout: it.TimeoutDuration(),
			States:  states,
		},
	)
}
