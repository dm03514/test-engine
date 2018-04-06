package fulfillment

import (
	"github.com/dm03514/test-engine/actions"
	"github.com/dm03514/test-engine/transcons"
	"reflect"
	"testing"
)

func TestRegistry_Load_Default_Noop(t *testing.T) {
	r, err := NewRegistry()
	var fType NoopFulfillment
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	m := make(map[string]interface{})
	m["type"] = "no registered type"
	name := "test"
	f, err := r.Load(
		m,
		name,
		actions.DummyAction{},
		transcons.Conditions{},
	)
	if err != nil {
		t.Log(err)
	}

	if reflect.TypeOf(f) != reflect.TypeOf(fType) {
		t.Fatalf("Expected type %s, received: %s", reflect.TypeOf(f), reflect.TypeOf(fType))
	}
}
