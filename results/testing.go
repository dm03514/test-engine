package results

type singleStringResult struct {
	V string
}

func (s singleStringResult) Error() error { return nil }
func (s singleStringResult) ValueOfProperty(property string) (Value, error) {
	return StringValue{V: s.V}, nil
}
