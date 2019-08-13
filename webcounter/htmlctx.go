package main

import (
	"fmt"
	"time"

	rediscounter "github.com/slavrd/go-redis-counter"
)

// htmlCounter represents the context needed to render the HTML counter page
type htmlCounter struct {
	counter       *rediscounter.RedisCounter
	CtrValue      int64
	Time          time.Time
	RedisConnInfo string
}

// IncrBy increases the underlying rediscounter.RedisCounter counter by a
// and sets the CtrValue and Time to the new values
func (htmlC *htmlCounter) IncrBy(a int64) error {
	rv, err := htmlC.counter.IncrBy(a)
	if err != nil {
		return err
	}
	htmlC.CtrValue = rv
	htmlC.Time = time.Now()
	return nil
}

// newHTMLCounter returns htmlCounter with CtrValue set to the current value of the created underlying rediscounter.RedisCounter
// and Time to the current time.
// redisConnInfo is set according to the input parameters.
func newHTMLCounter(rAddr, rPass, rKey string, rDB int) (*htmlCounter, error) {

	rc, err := rediscounter.NewCounter(rAddr, rPass, rKey, rDB)
	if err != nil {
		return &htmlCounter{}, err
	}

	rcValue, err := rc.Get()
	if err != nil {
		return &htmlCounter{}, err
	}

	htmlRC := &htmlCounter{
		counter:       rc,
		CtrValue:      rcValue,
		Time:          time.Now(),
		RedisConnInfo: fmt.Sprintf("redis @ %s db=%v key=%q", rAddr, rDB, rKey),
	}

	return htmlRC, nil
}
