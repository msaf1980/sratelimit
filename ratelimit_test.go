package sratelimit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/msaf1980/sratelimit"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	rate := 100 // per second
	reqDelay := time.Second / time.Duration(rate)
	delay := reqDelay / time.Duration(3)

	rl := sratelimit.New(rate)

	count := 5
	for i := 0; i < count; i++ {
		var d int64
		start := time.Now()
		time.Sleep(delay)
		rl.Take()
		stepDuration := time.Now().Sub(start).Microseconds()
		if i == 0 {
			d = delay.Microseconds()
		} else {
			d = reqDelay.Microseconds()
		}
		assert.Condition(
			t,
			func() bool { return stepDuration >= d-d/10 && stepDuration <= d+d/10 },
			fmt.Sprintf("Step %d: get delay %d us instead of %d us", i, stepDuration, reqDelay.Microseconds()),
		)
	}
}

func TestUnlimited(t *testing.T) {
	now := time.Now()
	rl := sratelimit.NewUnlimited()
	for i := 0; i < 1000; i++ {
		rl.Take()
	}
	assert.Condition(t, func() bool { return time.Now().Sub(now) < 1*time.Millisecond }, "no artificial delay")
}
