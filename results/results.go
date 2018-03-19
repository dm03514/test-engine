package results

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type NamedResult struct {
	Name string

	Result
}

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

func (o Override) Apply(rs Results, src string) (string, error) {
	REPLACE_ALL := -1
	r, err := rs.Get(o.FromState)
	if err != nil {
		return "", err
	}
	v, err := r.ValueOfProperty(o.UsingProperty)
	if err != nil {
		return "", err
	}
	replaced := strings.Replace(src, o.ToReplace, v.String(), REPLACE_ALL)
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
