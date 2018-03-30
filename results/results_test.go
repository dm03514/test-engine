package results

import "testing"

func TestOverride_Apply(t *testing.T) {
	rs := New(
		NamedResult{
			Name:   "test",
			Result: singleStringResult{V: "!!!"},
		},
	)

	o := Override{
		FromState:     "test",
		UsingProperty: "dummy",
		ToReplace:     "helllo",
	}
	s, err := o.Apply(*rs, "helllo")
	if err != nil {
		t.Error(err)
	}
	if s != "!!!" {
		t.Errorf("%q != %q, Expecting %q, Received %q", s, "!!!", "!!!", s)
	}
}
