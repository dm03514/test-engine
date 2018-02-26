package engine

import "time"

type DurationRecorder func(d time.Duration, err error)

func NoopDurationRecorder(d time.Duration, err error) {}
