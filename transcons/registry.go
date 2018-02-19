package transcons

import (
	"fmt"
)

type transConsLoader func(map[string]interface{}) (TransCon, error)

type TransConsRegistry struct {
	m map[string]transConsLoader
}

func (tcr TransConsRegistry) Load(tcm map[string]interface{}) (TransCon, error) {
	t := tcm["type"].(string)
	loaderFn, ok := tcr.m[t]
	if !ok {
		return nil, fmt.Errorf("Unable to parse transaction condition type type %s", t)
	}
	return loaderFn(tcm)
}

func NewTransConsRegistry() (TransConsRegistry, error) {
	return TransConsRegistry{
		m: map[string]transConsLoader{
			"assertions.IntEqual": NewIntEqualFromMap,
		},
	}, nil
}
