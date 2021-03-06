package actions

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

// SubprocessResult from suprocess execution
type SubprocessResult struct {
	Output     []byte
	Returncode int
	err        error
}

// Error unpacks error from the result
func (s SubprocessResult) Error() error {
	return s.err
}

// ValueOfProperty retrieves a result value based on a property of the result
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

// Subprocess can execute a shell command
type Subprocess struct {
	CommandName string `mapstructure:"command_name"`
	Args        []string
	Overrides   []results.Override
	Type        string
}

func (s Subprocess) applyOverrides(rs results.Results) (string, []string, error) {
	var err error
	cn := s.CommandName
	args := make([]string, len(s.Args))
	for _, o := range s.Overrides {
		cn, err = o.Apply(rs, cn)
		if err != nil {
			return "", nil, err
		}

		for i, arg := range s.Args {
			arg, err = o.Apply(rs, arg)
			if err != nil {
				return "", nil, err
			}
			args[i] = arg
		}
		s.Args = args
	}

	return cn, s.Args, nil
}

// Execute s a shell command in a subprocess
func (s Subprocess) Execute(ctx context.Context, rs results.Results) (results.Result, error) {
	var stdout, stderr bytes.Buffer

	cn, args, err := s.applyOverrides(rs)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"component":    s.Type,
		"command":      cn,
		"args":         args,
		"execution_id": ctx.Value(ids.Execution("execution_id")),
	}).Info("Execute()")

	cmd := exec.Command(cn, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	out := stdout.String() + stderr.String()

	log.WithFields(log.Fields{
		"component":    s.Type,
		"command":      cn,
		"args":         args,
		"output":       out,
		"error":        err,
		"execution_id": ctx.Value(ids.Execution("execution_id")),
	}).Info("CombinedOutput()")

	returncode, err := s.returnCode(err)
	if err != nil {
		return nil, err
	}

	return SubprocessResult{
		Output:     []byte(out),
		Returncode: returncode,
	}, nil
}

func (s Subprocess) returnCode(err error) (int, error) {

	if err == nil {
		return 0, nil
	}

	// TODO not sure the best way to handle this
	// It looks like there are non-portable ways to handle this
	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	// need to check if the error message is cross-platform safe
	switch err.Error() {
	case "exit status 1":
		return 1, nil
	}

	return -1, err
}

// NewSubprocessFromMap parses generic map into subprocess action
func NewSubprocessFromMap(m map[string]interface{}) (Action, error) {
	var sp Subprocess
	err := mapstructure.Decode(m, &sp)
	return sp, err
}
