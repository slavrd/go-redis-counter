package main

import "time"

// metrics represents the collected site usage data
// in the form of map where the keys are URL paths and
// values are times each path was called.
type metrics struct {
	Data          map[string]int
	Mutex         chan struct{}
	StartTime     time.Time
	RedisConnInfo string
}

// newMetrics creates a metrics instance with empty data.
// StartTime is set to time.Now() and RedisConnInfo to rinfo
func newMetrics(rinfo string) *metrics {
	return &metrics{
		Data:          make(map[string]int),
		Mutex:         make(chan struct{}, 1),
		StartTime:     time.Now(),
		RedisConnInfo: rinfo,
	}
}

// add will add data in a concurency sage manner by
// acquiring a lock and releasing it when done.
func (m *metrics) add(p string) {
	m.Mutex <- struct{}{}
	m.Data[p]++
	<-m.Mutex
}
