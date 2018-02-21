package actions

import (
	"fmt"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type SubprocessResult struct {
	Output     []byte
	Returncode int
	err        error
}

func (s SubprocessResult) Error() error {
	return s.err
}

func (s SubprocessResult) ValueOfProperty(property string) (results.Value, error) {
	switch property {
	case "returncode":
		return results.IntValue{V: s.Returncode}, nil
	case "output":
		return results.StringValue{V: string(s.Output)}, nil
	default:
		return nil, fmt.Errorf("No property %s in %+v", property, s)
	}
}

type Subprocess struct {
	CommandName string `mapstructure:"command_name"`
	Args        []string
}

func (s Subprocess) Execute() (results.Result, error) {
	cmd := exec.Command(s.CommandName, s.Args...)
	out, err := cmd.CombinedOutput()
	log.Infof("Execute() Subprocess out: %s, err: %s", out, err)
	if err != nil {
		return nil, err
	}

	return SubprocessResult{
		Output:     out,
		Returncode: 0,
	}, nil
}

func NewSubprocessFromMap(m map[string]interface{}) (Action, error) {
	var sp Subprocess
	err := mapstructure.Decode(m, &sp)
	return sp, err
}
