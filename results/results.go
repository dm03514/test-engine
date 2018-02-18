package results

type Result interface {
	Error() error
	ValueOfProperty(property string) (Value, error)
}

type Value interface {
	Int() int
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

type IntValue struct {
	V int
}

func (iv IntValue) Int() int {
	return iv.V
}
