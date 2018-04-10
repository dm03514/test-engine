// +build unit

package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
	"testing"
)

func TestSubprocess_applyOverrides_(t *testing.T) {
	t.Skip()
}

func TestSubprocessResult_ValueOfProperty(t *testing.T) {
	t.Skip()
}

func TestSubprocess_Execute(t *testing.T) {
	s := Subprocess{
		CommandName: "printf",
		Args:        []string{"hello"},
	}

	r, err := s.Execute(context.Background(), results.Results{})
	if err != nil {
		t.Error(err)
	}

	sr := r.(SubprocessResult)
	if string(sr.Output) != "hello" {
		t.Errorf("Expected %q received %q", "hello", sr.Output)
	}
	if sr.Returncode != 0 {
		t.Errorf("Expected returncode 0, received: %q", sr.Returncode)
	}

}
