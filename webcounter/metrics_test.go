package main

import (
	"testing"
	"time"
)

// TestNewMetrics checks if the returned metrics instance
// is initialized correctly.
func TestNewMetrics(t *testing.T) {
	m := newMetrics(redisConnInfo)

	if m.Data == nil {
		t.Error("Data property is nil")
	}

	if m.Mutex == nil {
		t.Error("Mutex property is nil")
	}

	if cap(m.Mutex) != 1 {
		t.Errorf("Mutex capacity is wrong, got:%v want: 1", cap(m.Mutex))
	}

	if m.RedisConnInfo != redisConnInfo {
		t.Errorf("RedisConnInfo property is wrong, got:%q want: %q", m.RedisConnInfo, redisConnInfo)
	}

	if m.StartTime == *new(time.Time) {
		t.Error("StartTime property is not set")
	}
}
