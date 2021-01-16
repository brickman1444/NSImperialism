package nationstates_api

import (
	"sync"
	"time"
)

type RateLimiter struct {
	queue            []time.Time
	numberOfRequests int
	perDuration      time.Duration
	mutex            sync.Mutex
}

func NewRateLimiter(numberOfRequests int, perDuration time.Duration) RateLimiter {
	return RateLimiter{
		queue:            make([]time.Time, 0, numberOfRequests),
		numberOfRequests: numberOfRequests,
		perDuration:      perDuration,
		mutex:            sync.Mutex{},
	}
}

func (limiter *RateLimiter) AddRequestTime(curerntTime time.Time) {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	limiter.queue = append(limiter.queue, curerntTime)
}

func (limiter *RateLimiter) IsAtRateLimit(currentTime time.Time) bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	if len(limiter.queue) < limiter.numberOfRequests {
		return false
	}

	howLongSinceMaxRequestNumber := currentTime.Sub(limiter.queue[limiter.numberOfRequests-1])

	return howLongSinceMaxRequestNumber < limiter.perDuration
}
