// +build integration

package actions

import (
	"context"
	"github.com/dm03514/test-engine/results"
	"testing"
)

func TestHTTP_Execute(t *testing.T) {
	h := HTTP{
		URL:    "https://google.com",
		Method: "GET",
	}

	r, err := h.Execute(context.Background(), results.Results{})
	if err != nil {
		t.Error(err)
	}

	result := r.(HTTPResult)
	if result.StatusCode != 200 {
		t.Errorf("Expected 200 received %q", result.StatusCode)
	}
}
