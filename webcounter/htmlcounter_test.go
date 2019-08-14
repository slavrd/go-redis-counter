package main

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func TestNewHTMLCounter(t *testing.T) {

	// set initial value for the counter in redis
	var initVal = int64(5)
	_, err := c.Set(rkey, initVal, 0).Result()
	if err != nil && err != redis.Nil {
		t.Fatalf("error setting initial counter value in redis: %v", err)
	}

	htmlC, err := newHTMLCounter(raddr, rpass, rkey, rdb)
	if err != nil {
		t.Fatalf("error createting webcounter.htmlCounter: %v", err)
	}

	if htmlC.counter == nil {
		t.Errorf("htmlCounter: counter was nil ")
	}

	if htmlC.CtrValue != initVal {
		t.Errorf("htmlCounter: unexpected initial CtrValue want: %v got: %v", initVal, htmlC.CtrValue)
	}

	var wantConnInfo = fmt.Sprintf("redis @ %s; db: %v; key: %q", raddr, rdb, rkey)
	if htmlC.RedisConnInfo != wantConnInfo {
		t.Errorf("htmlCounter: wrong RedisConnInfo, want: %q, got: %q", wantConnInfo, htmlC.RedisConnInfo)
	}

}
