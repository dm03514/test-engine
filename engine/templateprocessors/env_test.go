// +build unit

package templateprocessors

import (
	"reflect"
	"testing"
)

type StubEnv struct {
	m map[string]string
}

func (se StubEnv) LookupEnv(key string) (string, bool) {
	s, ok := se.m[key]
	return s, ok
}

func TestEnv_Process(t *testing.T) {
	testCases := []struct {
		name     string
		b        []byte
		expected []byte
		envvars  map[string]string
	}{
		{
			"pass_through",
			[]byte("hello"),
			[]byte("hello"),
			nil,
		},
		{
			"single_substitution",
			[]byte("hello $ENV_HI"),
			[]byte("hello you"),
			map[string]string{
				"HI": "you",
			},
		},
		{
			"single_substitution_multiple_targets",
			[]byte("hello $ENV_HI, $ENV_HI"),
			[]byte("hello you, you"),
			map[string]string{
				"HI": "you",
			},
		},
		{
			"multiple_substitutions",
			[]byte("hello $ENV_HI, $ENV_HOLA"),
			[]byte("hello you, friend"),
			map[string]string{
				"HI":   "you",
				"HOLA": "friend",
			},
		},
		{
			"multiple_substitutions_missing_one_target",
			[]byte("hello $ENV_HI, $ENV_HOLA"),
			[]byte("hello you, $ENV_HOLA"),
			map[string]string{
				"HI": "you",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ep := NewEnv(StubEnv{tc.envvars}.LookupEnv)
			out := ep.Process(tc.b)
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected %q Received %q", tc.expected, out)
			}
		})
	}
}
