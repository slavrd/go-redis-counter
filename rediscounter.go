// Package rediscounter provides a counter that stores its value in a redis server
package rediscounter

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

// RedisCounter represents a counter that stores its value in redis
type RedisCounter struct {
	rkey    string
	rclient *redis.Client
}

// Get returns the current value of the counter
func (rc *RedisCounter) Get() (int64, error) {

	result, err := rc.rclient.Get(rc.rkey).Result()

	// set rkey to 0 if it does not exist in redis
	if err != nil {
		return 0, fmt.Errorf("error reading form reids: %v", err)
	}

	// confirm that the value is int and return it
	cv, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("redis key: %q value: %q is not int: %v",
			rc.rkey, result, err)
	}
	return cv, nil
}

// IncrBy increases the counter's value by "a" amount and reports the resulting value
func (rc *RedisCounter) IncrBy(a int64) (int64, error) {
	rv, err := rc.rclient.IncrBy(rc.rkey, a).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing value: %v", err)
	}
	return rv, nil
}

// RedisHealth checks the redis server connection using PING
func (rc *RedisCounter) RedisHealth() error {
	_, err := rc.rclient.Ping().Result()
	return err
}

// NewCounter creates a RedisCounter with the provided connection details.
//
// "raddr" format should be host:port
//
// "rpass" can be set to "" if no password authentication is required.
func NewCounter(raddr, rpass, rkey string, rdb int) (*RedisCounter, error) {

	// create the redis client and check it can connect
	rclient := redis.NewClient(&redis.Options{
		Addr:     raddr,
		DB:       rdb,
		Password: rpass,
	})

	// Check if rkey exist and create it if it does not.
	_, err := rclient.Get(rkey).Result()
	if err == redis.Nil {
		_, err = rclient.Set(rkey, 0, 0).Result()
	}
	if err != nil {
		return nil, fmt.Errorf("failed connecting to redis: %v", err)
	}

	return &RedisCounter{rclient: rclient, rkey: rkey}, nil
}
