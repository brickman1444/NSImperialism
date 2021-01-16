package nationstates_api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterWithJustAFewRequestsWithinTimeLimitIsntAtLimit(t *testing.T) {

	duration, _ := time.ParseDuration("10m")
	limiter := NewRateLimiter(3, duration)

	tenTen, err := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)

	tenEleven, err := time.Parse(time.RFC3339, "2010-10-10T10:11:00Z")
	assert.NoError(t, err)

	assert.False(t, limiter.IsAtRateLimit(tenEleven))
}

func TestRateLimiterWithMaximumRequestsWithinTimeLimitIsAtLimit(t *testing.T) {

	duration, _ := time.ParseDuration("10m")
	limiter := NewRateLimiter(3, duration)

	tenTen, err := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)

	tenEleven, err := time.Parse(time.RFC3339, "2010-10-10T10:11:00Z")
	assert.NoError(t, err)

	assert.True(t, limiter.IsAtRateLimit(tenEleven))
}

func TestRateLimiterWithMaximumRequestsPastTimeLimitIsAtLimit(t *testing.T) {

	duration, _ := time.ParseDuration("10m")
	limiter := NewRateLimiter(3, duration)

	tenTen, err := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)

	tenThirty, err := time.Parse(time.RFC3339, "2010-10-10T10:30:00Z")
	assert.NoError(t, err)

	assert.False(t, limiter.IsAtRateLimit(tenThirty))
}

func TestRateLimiterWithHistoryAndMaximumRequestsWithinTimeLimitIsAtLimit(t *testing.T) {

	duration, _ := time.ParseDuration("10m")
	limiter := NewRateLimiter(3, duration)

	ten, err := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(ten)
	limiter.AddRequestTime(ten)
	limiter.AddRequestTime(ten)

	tenTen, err := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)

	tenEleven, err := time.Parse(time.RFC3339, "2010-10-10T10:11:00Z")
	assert.NoError(t, err)

	assert.True(t, limiter.IsAtRateLimit(tenEleven))
}

func TestRateLimiterWithHistoryAndMaximumRequestsPastTimeLimitIsAtLimit(t *testing.T) {

	duration, _ := time.ParseDuration("10m")
	limiter := NewRateLimiter(3, duration)

	ten, err := time.Parse(time.RFC3339, "2010-10-10T10:00:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(ten)
	limiter.AddRequestTime(ten)
	limiter.AddRequestTime(ten)

	tenTen, err := time.Parse(time.RFC3339, "2010-10-10T10:10:00Z")
	assert.NoError(t, err)

	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)
	limiter.AddRequestTime(tenTen)

	tenThirty, err := time.Parse(time.RFC3339, "2010-10-10T10:30:00Z")
	assert.NoError(t, err)

	assert.False(t, limiter.IsAtRateLimit(tenThirty))
}
