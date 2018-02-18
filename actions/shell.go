package actions

import (
	"fmt"
	"github.com/dm03514/test-engine/results"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type SubprocessResult struct {
	output     []byte
	returncode int
	err        error
}

func (s SubprocessResult) Error() error {
	return s.err
}

func (s SubprocessResult) ValueOfProperty(property string) (results.Value, error) {
	if property == "returncode" {
		return results.IntValue{V: s.returncode}, nil
	} else {
		return nil, fmt.Errorf("No property %s in %+v", property, s)
	}
}

type Subprocess struct {
	CommandName string
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
		output:     out,
		returncode: 0,
	}, nil
}

func NewSubprocess(a map[string]interface{}) (Action, error) {
	args := []string{}

	for _, arg := range a["args"].([]interface{}) {
		args = append(args, arg.(string))
	}

	return Subprocess{
		CommandName: a["command_name"].(string),
		Args:        args,
	}, nil
}
