package main

import (
	"fmt"
	"testing"

	rediscounter "github.com/slavrd/go-redis-counter"
)

// TestNewCounterCtx run tests by passing a mock function as argument
// and checks the values returned by NewCounterCtx.
func TestNewCounterCtx(t *testing.T) {
	tests := []struct {
		fResult int64
		fErr    error
	}{
		{
			fResult: int64(18),
			fErr:    nil,
		},
		{
			fResult: int64(18),
			fErr:    fmt.Errorf("this is not an error"),
		},
	}

	for _, test := range tests {

		f := func() (int64, error) {
			return test.fResult, test.fErr
		}
		r, err := newCounterCtx(f)

		if err != nil {

			if err.Error() != test.fErr.Error() {
				t.Errorf("returned wrong error value want: %v, got: %v", test.fErr.Error(), err.Error())
			}
			// if newCounterCtx returned an error, the result should be nil
			if r != nil {
				t.Errorf("counterCtx was not nil on error got: %v", r)
			}

		} else {

			if r.CtrValue != test.fResult {
				t.Errorf("CtrValue in counterCtx is wrong, want: %v, got: %v", test.fResult, r.CtrValue)
			}
		}
	}
}

// TestNewCounterWrapper tests the wrapping functions for NewCounterCtx.
// The testcase for each wrapper is defined as a struct and then the tests are executed over each testcase
func TestNewCounterWrappers(t *testing.T) {
	type testcaseWrapper struct {
		tFunc        func(*rediscounter.RedisCounter) (*counterCtx, error)
		name         string
		wantCountDif int64
		wantCtx      *counterCtx
	}

	tests := []testcaseWrapper{
		{
			tFunc:        newGetCtx,
			name:         "newGetCtx",
			wantCountDif: 0,
		},
		{
			tFunc:        newIncrCtx,
			name:         "newIncrCtx",
			wantCountDif: 1,
		},
	}

	for _, test := range tests {

		initVal, err := counter.Get()
		if err != nil {
			t.Fatalf("error getting initial counter value: %v", err)
		}

		r, err := test.tFunc(counter)
		if err != nil {
			t.Errorf("testcase: %s : func returned error: %v", test.name, err)
			continue
		}

		wv := initVal + test.wantCountDif

		if r.CtrValue != wv {
			t.Errorf("testcase: %s : func returned ctx with wrong count value: want: %v got: %v",
				test.name, wv, r.CtrValue)
		}
	}
}
