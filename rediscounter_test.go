/*Tests require a functioning redis instance with password authentication configured.
The redis address (host:port) must be set in $REDIS_ADDR variable.
If password authentication is set the password must be set in $REDIS_PASS variable.
The tests will be done with db 0 and key "count". If the key exists the test will delete it!
*/
package rediscounter

import (
	"log"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/go-redis/redis"
)

// redis key to use for the tests
var rkey = "count"

// a redis.Client to use for setting up redis test values
var c *redis.Client

func init() {

	// create the redis client
	c = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		DB:       0,
		Password: os.Getenv("REDIS_PASS"),
	})

	// confirm the client is working
	_, err := c.Ping().Result()
	if err != nil {
		log.Fatalf("error creating redis.Client: %v", err)
	}
}

// TestNewCounter confirms that a RedisCounter with a working redis.Client can be created
func TestNewCounter(t *testing.T) {

	// delete the "count" key if present in redis
	_, err := c.Del(rkey).Result()
	if err != nil && err != redis.Nil {
		t.Fatalf("failed to delete redis key %q: %v", rkey, err)
	}

	raddr := os.Getenv("REDIS_ADDR")
	rpass := os.Getenv("REDIS_PASS")
	rc, err := NewCounter(raddr, rpass, rkey, 0)
	if err != nil {
		t.Fatalf("error createting rediscounter.RedisCounter: %v", err)
	}

	if rc.rkey != rkey {
		t.Errorf("rediscounter.RedisCounter wrong key, want: %q, got: %q", rkey, rc.rkey)
	}

	rcOptions := rc.rclient.Options()
	if rcOptions.Addr != raddr {
		t.Errorf("rediscounter.RedisCounter wrong address, want: %q, got: %q", raddr, rcOptions.Addr)
	}
	if rcOptions.Password != rpass {
		t.Errorf("rediscounter.RedisCounter wrong password, want: %q, got: %q", rpass, rcOptions.Password)
	}
	if rcOptions.DB != 0 {
		t.Errorf("rediscounter.RedisCounter wrong database, want: 0, got: %q", rcOptions.DB)
	}
}

// TestGet verifies the RedisCounter.Get() method
func TestGet(t *testing.T) {

	// crete a rediscounter.RedisClient
	raddr := os.Getenv("REDIS_ADDR")
	rpass := os.Getenv("REDIS_PASS")
	rc, err := NewCounter(raddr, rpass, rkey, 0)
	if err != nil {
		t.Fatalf("error createting rediscounter.RedisCounter: %v", err)
	}

	// set the expected value in redis
	tests := []int64{0, 1, 7, math.MaxInt64}
	for _, tv := range tests {

		// set test value in redis
		_, err = c.Set(rkey, tv, 0).Result()
		if err != nil {
			t.Fatalf("error setting redis test value: %v", err)
		}

		rv, err := rc.Get()
		if err != nil {
			t.Fatalf("error retrieveing test value: %v", err)
		}
		if rv != tv {
			t.Fatalf("got wrong value, want: %v got: %v", tv, rv)
		}
	}
}

// TestIncrBy verifies the RedisCounter.IncrBy(a int64) method
func TestIncrBy(t *testing.T) {

	// set the redis key value to 0
	_, err := c.Set(rkey, 0, 0).Result()
	if err != nil {
		t.Fatalf("error setting redis test value: %v", err)
	}

	// crete a rediscounter.RedisClient
	raddr := os.Getenv("REDIS_ADDR")
	rpass := os.Getenv("REDIS_PASS")
	rc, err := NewCounter(raddr, rpass, rkey, 0)
	if err != nil {
		t.Fatalf("error createting rediscounter.RedisCounter: %v", err)
	}

	for i := int64(0); i < 10; i++ {
		// get current value directly from redis
		rv, err := c.Get(rkey).Result()
		if err != nil {
			t.Errorf("error retrieving redis value before iteration %v: %v", i, err)
		}
		// convert the value to int64
		bIncr, _ := strconv.ParseInt(rv, 10, 64) // ignoring conversion error

		r, err := rc.IncrBy(i)
		if err != nil {
			t.Errorf("error increasing counter value by %v: %v", i, err)
		}

		// get the increased value from redis
		rv, err = c.Get(rkey).Result()
		if err != nil {
			t.Errorf("error retrieving redis value before iteration %v: %v", i, err)
		}
		// convert the value to int64
		aIncr, _ := strconv.ParseInt(rv, 10, 64) // ignoring conversion error

		if aIncr != bIncr+i {
			t.Errorf("IncrBy(%v) did not increase the counter with correct value. want: %v got: %v", i, bIncr+i, aIncr)
		}

		if r != aIncr {
			t.Errorf("IncrBy(%v) did not return actual redis value want: %v got: %v", i, aIncr, r)
		}

	}
}
