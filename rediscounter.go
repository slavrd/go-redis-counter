// Package rediscounter provides a counter that stores its value in a redis server
package rediscounter

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

// RedisCounter represents a counter that stores its value in redis
type RedisCounter struct {
	rkey    string
	rclient *redis.Client
	mutex   *sync.Mutex
}

// Get returns the current value of the counter
func (rc *RedisCounter) Get() (int64, error) {

	result, err := rc.rclient.Get(rc.rkey).Result()

	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("error reading form redis: %v", err)
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
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rv, err := rc.rclient.IncrBy(rc.rkey, a).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing value: %v", err)
	}
	return rv, nil
}

// DecrBy decreases the counter's value by "a" amount and reports the resulting value.
// The counter's value cannot go below 0.
func (rc *RedisCounter) DecrBy(a int64) (int64, error) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	v, err := rc.Get()
	if err != nil {
		return 0, err
	}

	// check if value is already 0 and simply return 0
	if v == 0 {
		return 0, nil
	}

	// check if value will go below 0 and if so set it to 0
	if v-a < 0 {
		_, err = rc.rclient.Set(rc.rkey, 0, 0).Result()
		if err != nil {
			return 0, fmt.Errorf("error decreasing counter value: %v", err)
		}
		return 0, nil
	}

	// decrease the counter value by a
	r, err := rc.rclient.DecrBy(rc.rkey, a).Result()
	if err != nil {
		return 0, fmt.Errorf("error decreasing counter value: %v", err)
	}
	return r, nil

}

// Reset sets the counter value to 0.
// It will allways return 0 so use error to determine if successful.
func (rc *RedisCounter) Reset() (int64, error) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	_, err := rc.rclient.Set(rc.rkey, 0, 0).Result()
	if err != nil {
		return 0, err
	}
	return 0, nil
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
		err = fmt.Errorf("failed connecting to redis: %v", err)
	}

	return &RedisCounter{rclient: rclient, rkey: rkey, mutex: &sync.Mutex{}}, err
}
