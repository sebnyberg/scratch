package main

import (
	"sync/atomic"
	"time"
)

// Courtesy of Egon Elbre
//
type RateLimiter struct {
	duration int64
	limit    int64

	deadline int64 // atomic
	writes   int64 // atomic
}

func NewLimiter(duration time.Duration, limit int64) *RateLimiter {
	return &RateLimiter{
		duration: duration.Nanoseconds(),
		limit:    limit,
	}
}

func (lim *RateLimiter) Exceeded() bool {
	deadline := atomic.LoadInt64(&lim.deadline)
	t := time.Now().UnixNano()

	if deadline < t {
		atomic.StoreInt64(&lim.writes, 0)
		atomic.StoreInt64(&lim.deadline, t+lim.duration)
		return false
	}

	count := atomic.AddInt64(&lim.writes, 1)
	if count > lim.limit {
		return true
	}
	return false
}
