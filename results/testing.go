package results

type SingleStringResult struct {
	V string
}

func (s SingleStringResult) Error() error { return nil }
func (s SingleStringResult) ValueOfProperty(property string) (Value, error) {
	return StringValue{V: s.V}, nil
}
