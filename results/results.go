package results

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// NamedResult has a result with string identifier
type NamedResult struct {
	Name string

	Result
}

// Result interface can identify errors, or fetch values based on string identifiers
type Result interface {
	Error() error
	ValueOfProperty(property string) (Value, error)
}

// Value can return an Int() or String() value
type Value interface {
	Int() int
	String() string
}

// ErrorResult  contains an error and reference to another result
type ErrorResult struct {
	From Result
	Err  error
}

// Error returns the erro that is wrapped
func (er ErrorResult) Error() error {
	return er.Err
}

// ValueOfProperty is a dummy to fulfill Result interface
func (er ErrorResult) ValueOfProperty(property string) (Value, error) {
	return nil, nil
}

// DummyIntValue contains nothing
type DummyIntValue struct{}

// Int returns 0
func (div DummyIntValue) Int() int { return 0 }

// IntValue represents an int
type IntValue struct {
	V int
}

// Int returns the int value that was initialized
func (iv IntValue) Int() int {
	return iv.V
}

// String returns the string representation of the original int
func (iv IntValue) String() string {
	return strconv.FormatInt(int64(iv.V), 10)
}

// StringValue contains a string
type StringValue struct {
	DummyIntValue

	V string
}

// String returns the initialized string
func (sv StringValue) String() string {
	return sv.V
}

// Override contains result
type Override struct {
	FromState     string `mapstructure:"from_state"`
	UsingProperty string `mapstructure:"using_property"`
	ToReplace     string `mapstructure:"to_replace"`
}

func (o Override) Apply(rs Results, src string) (string, error) {
	replaceAll := -1
	r, err := rs.Get(o.FromState)
	if err != nil {
		return "", err
	}
	v, err := r.ValueOfProperty(o.UsingProperty)
	if err != nil {
		return "", err
	}
	replaced := strings.Replace(src, o.ToReplace, v.String(), replaceAll)
	log.WithFields(log.Fields{
		"component": "results.Override",
		"replacing": o.ToReplace,
		"with":      v.String(),
		"in":        src,
		"result":    replaced,
	}).Info("apply()")
	return replaced, nil
}

type Results struct {
	byStateName map[string]Result
}

func (rs *Results) Add(r NamedResult) {
	log.WithFields(log.Fields{
		"component": "results",
		"name":      r.Name,
		"adding":    r.Name,
	}).Info("Add()")
	rs.byStateName[r.Name] = r
}

func (rs *Results) Get(k string) (Result, error) {
	r, ok := rs.byStateName[k]
	log.WithFields(log.Fields{
		"component": "results",
		"name":      k,
	}).Info("Get()")

	if !ok {
		return nil, fmt.Errorf("Unable to find key: %s in %+v", k, rs)
	}
	return r, nil
}

func New(nrs ...NamedResult) *Results {
	rs := &Results{
		byStateName: make(map[string]Result),
	}
	for _, r := range nrs {
		rs.Add(r)
	}
	return rs
}
