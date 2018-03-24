package engine

import "time"

// StateDurationRecorder for recording a duration on a state
type StateDurationRecorder func(stateName string, testName string, d time.Duration, err error)

// TestDurationRecorder for recording duration on a test execution
type TestDurationRecorder func(testName string, d time.Duration, err error)

// NoopStateDurationRecorder doesn't do anything! it's the default
func NoopStateDurationRecorder(stateName string, testName string, d time.Duration, err error) {}

// NoopTestDurationRecorder doesn't do anything! it's just the default
func NoopTestDurationRecorder(testName string, d time.Duration, err error) {}
