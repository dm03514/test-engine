package actions

import "github.com/dm03514/test-engine/results"

type Action interface {
	Execute(results.Results) (results.Result, error)
}
