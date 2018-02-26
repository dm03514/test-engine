package engine

type Option func(e *Engine)

func OptionRecordStateDuration(f DurationRecorder) func(e *Engine) {
	return func(e *Engine) {
		e.recordStateDuration = f
	}
}

func OptionRecordTestDuration(f DurationRecorder) func(e *Engine) {
	return func(e *Engine) {
		e.recordTestDuration = f
	}
}
