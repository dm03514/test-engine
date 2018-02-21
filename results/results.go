package results

type Result interface {
	Error() error
	ValueOfProperty(property string) (Value, error)
}

type Value interface {
	Int() int
	String() string
}

type ErrorResult struct {
	From Result
	Err  error
}

func (er ErrorResult) Error() error {
	return er.Err
}
func (er ErrorResult) ValueOfProperty(property string) (Value, error) {
	return nil, nil
}

type DummyStringValue struct{}

func (dsv DummyStringValue) String() string { return "" }

type DummyIntValue struct{}

func (div DummyIntValue) Int() int { return 0 }

type IntValue struct {
	DummyStringValue

	V int
}

func (iv IntValue) Int() int {
	return iv.V
}

type StringValue struct {
	DummyIntValue

	V string
}

func (sv StringValue) String() string {
	return sv.V
}
