package mocks

import (
	"time"
)

type FixedClock struct {
	now time.Time
}

func NewFixedClock(now time.Time) *FixedClock {
	return &FixedClock{now: now.UTC()}
}

func (c *FixedClock) Now() time.Time {
	return c.now
}

func (c *FixedClock) Today() (int, time.Month, int) {
	return c.now.Year(), c.now.Month(), c.now.Day()
}
