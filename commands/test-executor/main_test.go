// +build functional

package main

import (
	"flag"
	"path"
	"testing"
)

var rootTestDir string

func init() {
	flag.StringVar(&rootTestDir, "root-test-dir", "", "full path to test dir")
	flag.Parse()
}

func TestMain(t *testing.T) {
	testCases := []struct {
		testFile string
	}{
		{"subprocess_exit_code.yml"},
		{"subprocess_multiple_conditions.yml"},
		{"multiple_states.yml"},
		{"previous_state_overrides.yml"},
	}
	for _, tc := range testCases {
		t.Run(tc.testFile, func(t *testing.T) {
			err := executeTest(path.Join(rootTestDir, tc.testFile))
			if err != nil {
				t.Error(err)
			}
		})
	}
}
