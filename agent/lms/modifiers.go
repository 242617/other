package lms

import "time"

type Modifier = func(*LMS)

func WithHost(host string) Modifier { return func(lms *LMS) { lms.host = host } }

func withDefaultClientTimeout() Modifier         { return WithClientTimeout(time.Minute) }
func WithClientTimeout(d time.Duration) Modifier { return func(lms *LMS) { lms.timeout = d } }
