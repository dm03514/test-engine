package engine

import "time"

type StateDurationRecorder func(stateName string, testName string, d time.Duration, err error)
type TestDurationRecorder func(testName string, d time.Duration, err error)

func NoopStateDurationRecorder(stateName string, testName string, d time.Duration, err error) {}

func NoopTestDurationRecorder(testName string, d time.Duration, err error) {}
