package engine

// Option function allows configuration of an engine instance
type Option func(e *Engine)

// OptionRecordStateDuration allow setting a function to be called to record state durations
func OptionRecordStateDuration(f StateDurationRecorder) func(e *Engine) {
	return func(e *Engine) {
		e.recordStateDuration = f
	}
}

// OptionRecordTestDuration allow setting a function to be called to record test durations
func OptionRecordTestDuration(f TestDurationRecorder) func(e *Engine) {
	return func(e *Engine) {
		e.recordTestDuration = f
	}
}
