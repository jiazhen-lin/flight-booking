package service

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Allow(ctx context.Context) (bool, error)
}

var (
	keyToTokenBucketLimiter = sync.Map{}
)

// NOTE: this is a simple implementation of token bucket rate limiter implemented by golang.org/x/time/rate
// We can use redis to implement for concurrency
type tokenBucketRateLimiter struct {
	duration time.Duration
	burst    int
	limiter  *rate.Limiter
}

var _ RateLimiter = (*tokenBucketRateLimiter)(nil)

type TokenBucketConfig struct {
	Key      string
	Duration time.Duration
	Burst    int
}

func NewMemoryTokenBucketRateLimiter(conf TokenBucketConfig) (*tokenBucketRateLimiter, error) {
	v, _ := keyToTokenBucketLimiter.LoadOrStore(conf.Key,
		&tokenBucketRateLimiter{
			duration: conf.Duration,
			burst:    conf.Burst,
			limiter:  rate.NewLimiter(rate.Every(conf.Duration), conf.Burst),
		},
	)
	limiter, ok := v.(*tokenBucketRateLimiter)
	if !ok {
		return nil, fmt.Errorf("convert value to rate.Limiter error, got: %v", reflect.TypeOf(v).Kind())
	}
	return limiter, nil
}

func (r *tokenBucketRateLimiter) Allow(ctx context.Context) (bool, error) {
	return r.limiter.Allow(), nil
}
