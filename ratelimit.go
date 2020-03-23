package sratelimit

import (
	"time"
)

// Limiter is used to rate-limit some process, possibly across goroutines.
// The process is expected to call Take() before every iteration, which
// may block to throttle the goroutine.
type Limiter interface {
	// Take should block to make sure that the RPS is met.
	Take() time.Time
	TakeWithTime(now time.Time) time.Time
}

type limiter struct {
	last time.Time

	perRequest time.Duration
}

// New returns a Limiter that will limit to the given RPS.
func New(rate int) Limiter {
	l := &limiter{
		time.Time{},
		time.Second / time.Duration(rate),
	}

	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *limiter) TakeWithTime(now time.Time) time.Time {
	sleepFor := t.perRequest - now.Sub(t.last)

	if sleepFor > time.Duration(0) {
		time.Sleep(sleepFor)
		t.last = now.Add(sleepFor)
	} else {
		t.last = now
	}

	return now
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *limiter) Take() time.Time {
	return t.TakeWithTime(time.Now())
}

type unlimited struct{}

// NewUnlimited returns a RateLimiter that is not limited.
func NewUnlimited() Limiter {
	return unlimited{}
}

func (unlimited) Take() time.Time {
	return time.Now()
}

func (unlimited) TakeWithTime(now time.Time) time.Time {
	return time.Now()
}
