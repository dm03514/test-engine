package actions

import (
	"github.com/dm03514/test-engine/results"
	"net/http"
	"testing"
)

func TestHTTPResult_ValueOfProperty(t *testing.T) {
	testCases := []struct {
		property string
		err      bool
		value    results.Value
		resp     *http.Response
		body     []byte
	}{
		{"status_code", false, results.IntValue{}, &http.Response{}, nil},
		{"body", false, results.StringValue{}, &http.Response{}, nil},
		{"unknown", true, nil, &http.Response{}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.property, func(t *testing.T) {
			r := HTTPResult{
				Response: tc.resp,
				body:     tc.body,
			}
			v, err := r.ValueOfProperty(tc.property)
			// if this test is expecting an error than we need to error if the
			// error is empty!
			if tc.err && err == nil {
				t.Error(err)
			}
			if v != tc.value {
				t.Errorf("Expected %+v, received %+v", tc.value, v)
			}
		})
	}
}
