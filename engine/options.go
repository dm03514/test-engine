package engine

type Option func(e *Engine)

func OptionRecordStateDuration(f StateDurationRecorder) func(e *Engine) {
	return func(e *Engine) {
		e.recordStateDuration = f
	}
}

func OptionRecordTestDuration(f TestDurationRecorder) func(e *Engine) {
	return func(e *Engine) {
		e.recordTestDuration = f
	}
}
