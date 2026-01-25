package lms

import "time"

type Modifier = func(*LMS)

// "http://0.0.0.0:1234"
func WithAddress(address string) Modifier { return func(lms *LMS) { lms.address = address } }

func withDefaultClientTimeout() Modifier         { return WithClientTimeout(time.Minute) }
func WithClientTimeout(d time.Duration) Modifier { return func(lms *LMS) { lms.timeout = d } }
