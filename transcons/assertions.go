package transcons

import (
	"fmt"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type IntEqual struct {
	UsingProperty string `mapstructure:"using_property"`
	ToEqual       int    `mapstructure:"to_equal"`
}

func (ie IntEqual) Evaluate(r results.Result) results.Result {
	v, err := r.ValueOfProperty(ie.UsingProperty)
	log.Infof("IntEqual.Evaluate() result: %+v, against: %+v", r, v)
	if err != nil {
		return results.ErrorResult{
			From: r,
			Err:  err,
		}
	}

	if v.Int() != ie.ToEqual {
		return results.ErrorResult{
			From: r,
			Err:  fmt.Errorf("%+v != %+v", ie.ToEqual, v.Int()),
		}
	}
	return r
}

func NewIntEqualFromMap(m map[string]interface{}) (TransCon, error) {
	var ie IntEqual
	err := mapstructure.Decode(m, &ie)
	return ie, err
}

type StringEqual struct {
	UsingProperty string `mapstructure:"using_property"`
	ToEqual       string `mapstructure:"to_equal"`
}

func (se StringEqual) Evaluate(r results.Result) results.Result {
	v, err := r.ValueOfProperty(se.UsingProperty)
	log.Infof("StringEqual.Evaluate(%s) result: %+v, against: %+v.  Error: %+v",
		se.UsingProperty, r, v, err)
	if err != nil {
		return results.ErrorResult{
			From: r,
			Err:  err,
		}
	}
	if v.String() != se.ToEqual {
		return results.ErrorResult{
			From: r,
			Err: fmt.Errorf("%q != %q, expected %q, received %q",
				se.ToEqual, v.String(), se.ToEqual, v.String()),
		}
	}
	return r
}

func NewStringEqualFromMap(m map[string]interface{}) (TransCon, error) {
	var se StringEqual
	err := mapstructure.Decode(m, &se)
	return se, err
}
