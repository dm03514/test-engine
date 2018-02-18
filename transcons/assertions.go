package transcons

import (
	"fmt"
	"github.com/dm03514/test-engine/results"
)

type IntEqual struct {
	UsingProperty string `mapstructure:"using_property"`
	ToEqual       int    `mapstructure:"to_equal"`
}

func (ie IntEqual) Evaluate(r results.Result) results.Result {
	v, err := r.ValueOfProperty(ie.UsingProperty)
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
