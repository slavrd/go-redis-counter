package main

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	rediscounter "github.com/slavrd/go-redis-counter"
)

// TestNewHandler passes a mock function that returns a known *counterCtx.
// It then calls the resulting http.Handler and checks the response it writes.
func TestNewHandler(t *testing.T) {

	type testcaseNewHandler struct {
		ctxf     func(*rediscounter.RedisCounter) (*counterCtx, error)
		wantCode int
		name     string
		wantBody []byte
	}

	testCtx := &counterCtx{
		CtrValue:      5,
		RedisConnInfo: redisConnInfo,
		Time:          time.Now(),
	}

	okBodyBuf := bytes.NewBuffer(make([]byte, 0))
	err := htmlCounterTpl.Execute(okBodyBuf, testCtx)
	if err != nil {
		t.Fatalf("error setting up response body check: %v", err)
	}

	tests := []testcaseNewHandler{
		{
			name: "ok",
			ctxf: func(c *rediscounter.RedisCounter) (*counterCtx, error) {
				if !reflect.DeepEqual(c, counter) {
					t.Error("argument of ctxf is not the global counter")
				}
				return testCtx, nil
			},
			wantCode: 200,
			wantBody: okBodyBuf.Bytes(),
		},
		{
			wantCode: 500,
			name:     "error",
			ctxf: func(c *rediscounter.RedisCounter) (*counterCtx, error) {
				return nil, fmt.Errorf("this is not an error")
			},
			wantBody: []byte("Internal server error!"),
		},
	}

	for _, test := range tests {

		h := newHandler(test.ctxf, htmlCounterTpl)

		r := httptest.NewRequest("GET", "/testpath", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("testcase %q: wrong status code, want: %v, got: %v", test.name, test.wantCode, w.Code)
		}

		if !bytes.Equal(w.Body.Bytes(), test.wantBody) {
			t.Errorf("handler testcase: %v\n wrong body \nwant:\n\n%s\n\ngot:\n\n%s\n", test.name, string(test.wantBody), w.Body.String())
		}
	}
}
