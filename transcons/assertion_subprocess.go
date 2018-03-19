package transcons

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type Subprocess struct {
	UsingProperty string `mapstructure:"using_property"`
	ToEqual       string `mapstructure:"to_equal"`

	CommandName string `mapstructure:"command_name"`
	Args        []string
	Type        string
}

func (s Subprocess) substituteProperty(v results.Value) (string, []string, error) {
	REPLACE_ALL := -1
	toReplace := fmt.Sprintf("$%s", s.UsingProperty)
	command := strings.Replace(s.CommandName, toReplace, v.String(), REPLACE_ALL)

	args := make([]string, len(s.Args))
	for i, arg := range s.Args {
		args[i] = strings.Replace(arg, toReplace, v.String(), REPLACE_ALL)
	}
	return command, args, nil
}

func (s Subprocess) Evaluate(ctx context.Context, r results.Result) results.Result {
	v, err := r.ValueOfProperty(s.UsingProperty)
	log.WithFields(log.Fields{
		"component":      s.Type,
		"execution_id":   ctx.Value(ids.Execution("execution_id")),
		"using_property": s.UsingProperty,
		"to_equal":       s.ToEqual,
		"against":        v.String(),
	}).Info("Evaluate() result: %+v, against: %+v", r, v)
	if err != nil {
		return results.ErrorResult{
			From: r,
			Err:  err,
		}
	}

	cn, args, err := s.substituteProperty(v)
	if err != nil {
		return results.ErrorResult{
			From: r,
			Err:  err,
		}
	}

	log.Infof("Subprocess.Evaluate() command: `%s` args: `%s`", cn, args)
	cmd := exec.Command(cn, args...)
	out, err := cmd.CombinedOutput()
	log.Infof("Subprocess.Evaluate() Subprocess out: %q, err: %s", out, err)
	if err != nil {
		return results.ErrorResult{
			From: r,
			Err:  err,
		}
	}

	if s.ToEqual != string(out) {
		return results.ErrorResult{
			From: r,
			Err: fmt.Errorf("%q != %q, expected %q, received %q",
				s.ToEqual, out, s.ToEqual, string(out)),
		}
	}
	return r

}

func NewSubprocessFromMap(m map[string]interface{}) (TransCon, error) {
	var s Subprocess
	err := mapstructure.Decode(m, &s)
	return s, err
}
