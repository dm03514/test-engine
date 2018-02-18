package transcons

import "github.com/dm03514/test-engine/actions"

type TransCon interface {
	Evaluate(result actions.Result) actions.Result
}
