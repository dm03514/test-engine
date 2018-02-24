package results

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

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

type DummyIntValue struct{}

func (div DummyIntValue) Int() int { return 0 }

type IntValue struct {
	V int
}

func (iv IntValue) Int() int {
	return iv.V
}

func (iv IntValue) String() string {
	return strconv.FormatInt(int64(iv.V), 10)
}

type StringValue struct {
	DummyIntValue

	V string
}

func (sv StringValue) String() string {
	return sv.V
}

type Override struct {
	FromState     string `mapstructure:"from_state"`
	UsingProperty string `mapstructure:"using_property"`
	ToReplace     string `mapstructure:"to_replace"`
}

func (o Override) Apply(rs Results, target string) (string, error) {
	REPLACE_ALL := -1
	r, err := rs.Get(o.FromState)
	if err != nil {
		return "", err
	}
	v, err := r.ValueOfProperty(o.UsingProperty)
	if err != nil {
		return "", err
	}
	replaced := strings.Replace(target, o.ToReplace, v.String(), REPLACE_ALL)
	log.Infof("Override.Apply() -> replacing %s with %s in %s -> %s",
		o.ToReplace, v.String(), target, replaced)
	return replaced, nil
}

type Results struct {
	byStateName map[string]Result
}

func (rs *Results) Add(k string, r Result) {
	log.Infof("Results.Add(%s, %+v)", k, r)
	rs.byStateName[k] = r
}

func (rs *Results) Get(k string) (Result, error) {
	r, ok := rs.byStateName[k]
	log.Infof("Results.Get(%s) -> %+v", k, r)
	if !ok {
		return nil, fmt.Errorf("Unable to find key: %s in %+v", k, rs)
	}
	return r, nil
}

func New() *Results {
	return &Results{
		byStateName: make(map[string]Result),
	}
}
