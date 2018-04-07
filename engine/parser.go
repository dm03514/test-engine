package engine

import (
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/engine/templateprocessors"
	"github.com/dm03514/test-engine/fulfillment"
	"github.com/dm03514/test-engine/transcons"
	"github.com/go-yaml/yaml"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type actionRegistry interface {
	Load(am map[string]interface{}) (actions.Action, error)
}

type transConsRegistry interface {
	Load(tcm map[string]interface{}) (transcons.TransCon, error)
}

type intermediaryState struct {
	Name                 string
	Fulfillment          map[string]interface{}
	Action               map[string]interface{}
	TransitionConditions []map[string]interface{} `yaml:"transition_conditions"`
}

func (is intermediaryState) ParsedTransCons(tcr transConsRegistry) (transcons.Conditions, error) {
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

func (is intermediaryState) State(ar actionRegistry, tcr transConsRegistry) (State, error) {
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

// NewFromYaml can parse a slice of bytes, as yaml, into a test!
func NewFromYaml(b []byte, ar actionRegistry, tcr transConsRegistry, f Factory) (*Engine, error) {
	it := intermediaryTest{}
	ep := templateprocessors.NewEnv(os.LookupEnv)
	uuidProcessor := templateprocessors.NewUUID(uuid.NewV4)

	err := yaml.Unmarshal(
		ep.Process(
			uuidProcessor.Process(
				b,
			),
		), &it,
	)
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
			"raw_state": ps.Name,
		}).Debug("parsing_tate")

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
