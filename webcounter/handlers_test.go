package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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

// TestNewHealthHandler uses a mock function to test the http.handlers returned by newHealthHandler
func TestNewHealthHandler(t *testing.T) {

	type testcaseNewHealth struct {
		name     string
		wantCode int
		wantBody []byte
		hcf      func() error
	}

	tests := []testcaseNewHealth{
		{
			name:     "no error",
			wantCode: 200,
			wantBody: []byte("OK"),
			hcf:      func() error { return nil },
		},
		{
			name:     "error",
			wantCode: 500,
			wantBody: []byte("Redis server is down!"),
			hcf:      func() error { return fmt.Errorf("this is not an error") },
		},
	}

	for _, test := range tests {
		h := newHealthHandler(test.hcf)

		r := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("testcase: %q returned wrong status code want: %v got: %v", test.name, test.wantCode, w.Code)
		}

		if !bytes.Equal(w.Body.Bytes(), test.wantBody) {
			t.Errorf("testcase: %q returned wrong body\nwant: %q\ngot: %q", test.name, test.wantBody, w.Body.Bytes())
		}
	}

}

// TestNewMetricsHandler will test the returned handler
// It will not use the global htmlMetricsTpl as it needs to change the template's timenow func to return a fixed value
func TestNewMetricsHandler(t *testing.T) {

	// create the metrics instance to test against
	// NOTE: calling the handler should change the Data
	tm := newMetrics(redisConnInfo)
	tm.Data["get"] = 5
	tm.Data["incr"] = 7

	// load the template file
	content, err := ioutil.ReadFile(*mtplPath)
	if err != nil {
		t.Fatalf("error reading template file: %v", err)
	}

	// create the html template with "timenow" func set to return static value
	tt := time.Now()
	tpl, err := template.New("metrics").
		Funcs(template.FuncMap{"timenow": func() time.Time { return tt }}).
		Parse(string(content))
	if err != nil {
		log.Fatalf("error loading metrics html template: %v", err)
	}

	// create the handler, invoke it, record the request
	h := newMetricsHandler(tm, tpl)
	rp := "/metrics"
	mv := tm.Data[rp]
	r := httptest.NewRequest("GET", rp, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	// confirm that the handler increased the value for the request path
	if tm.Data[rp] != mv+1 {
		t.Errorf("invoking the handler did not record the call in the metrics data\n, wanted: %v got: %v", mv+1, tm.Data[rp])
	}

	// define the expected request body
	// should be done after calling the handler so that it will reflect the data change
	buf := bytes.NewBuffer(make([]byte, 0))
	tpl.Execute(buf, tm) // ignorring potential error

	if w.Code != 200 {
		t.Errorf("wrong status code, wanted: 200 got: %v", w.Code)
	}

	if !bytes.Equal(buf.Bytes(), w.Body.Bytes()) {
		t.Errorf("wrong body \nwant:\n\n%s\n\ngot:\n\n%s\n", string(buf.Bytes()), w.Body.String())
	}
}

// TestNewCrashHandler will test the handler returned by newCrashHandler
// by checking if the argument function is called correctly.
func TestNewCrashHandler(t *testing.T) {

	var res interface{}
	f := func(a ...interface{}) {
		res = a[0]
	}

	want := "test message"
	h := newCrashHandler(f, want)
	r := httptest.NewRequest("GET", "/crash", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	got, ok := res.(string)

	if !ok {
		t.Fatalf("result variable is wrong type, got: %v want: string", reflect.TypeOf(res))
	}

	if got != want {
		t.Errorf("result variable holds wrong value, got: %q want: %q", got, want)
	}

}
