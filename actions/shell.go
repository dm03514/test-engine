package actions

import (
	"os/exec"
)

type SubprocessResult struct {
	output     []byte
	returncode int
}

type Subprocess struct {
	CommandName string
	Args        []string
}

func (s *Subprocess) Execute() (Result, error) {
	cmd := exec.Command(s.CommandName, s.Args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return SubprocessResult{
		output:     out,
		returncode: 0,
	}, nil
}
