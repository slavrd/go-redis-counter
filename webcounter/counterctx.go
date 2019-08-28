package main

import (
	"time"

	rediscounter "github.com/slavrd/go-redis-counter"
)

// counterCtx represents the counter data with which the html template is rendered for response
type counterCtx struct {
	CtrValue      int64
	Time          time.Time
	RedisConnInfo string
}

// newCounterCtx creates a *counterCtx instance with
// CtrValue set to the result of cf() and Time set to time.Now().
func newCounterCtx(cf func() (int64, error)) (*counterCtx, error) {
	r := &counterCtx{
		RedisConnInfo: redisConnInfo,
	}
	var err error
	r.CtrValue, err = cf()
	if err != nil {
		return nil, err
	}
	r.Time = time.Now()
	return r, nil
}

// Below are wrappers for the newCounterCtx that call it with specific RedisCounter methods
func newGetCtx(c *rediscounter.RedisCounter) (*counterCtx, error) {
	return newCounterCtx(c.Get)
}

func newIncrCtx(c *rediscounter.RedisCounter) (*counterCtx, error) {
	return newCounterCtx(func() (int64, error) { return c.IncrBy(1) })
}

func newDecrCtx(c *rediscounter.RedisCounter) (*counterCtx, error) {
	return newCounterCtx(func() (int64, error) { return c.DecrBy(1) })
}

func newResetCtx(c *rediscounter.RedisCounter) (*counterCtx, error) {
	return newCounterCtx(c.Reset)
}
