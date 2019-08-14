package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlers runs testcases over the server handlers.
func TestHandlers(t *testing.T) {

	// define testcases
	tests := []struct {
		handler  func(w http.ResponseWriter, r *http.Request)
		reqPath  string
		wantCode int
		ctrDif   int64
	}{
		{
			handler:  handleGet,
			reqPath:  "/get",
			wantCode: 200,
			ctrDif:   0,
		},
		{
			handler:  handleIncr,
			reqPath:  "/incr",
			wantCode: 200,
			ctrDif:   1,
		},
	}

	// execute tests
	for _, test := range tests {
		initCtrValue := htmlCounterCtx.CtrValue
		initTime := htmlCounterCtx.Time

		r := httptest.NewRequest("GET", test.reqPath, nil)
		w := httptest.NewRecorder()

		test.handler(w, r)

		if w.Code != test.wantCode {
			t.Errorf("handler testcase:\n%v\nwrong status code, want: %v, got: %v", test, test.wantCode, w.Code)
		}

		if initCtrValue+test.ctrDif != htmlCounterCtx.CtrValue {
			t.Errorf("handler testcase:\n%v\nincreased the counter value. Initial value: %v current value: %v", test, initCtrValue, htmlCounterCtx.CtrValue)
		}

		if htmlCounterCtx.Time.Sub(initTime) <= 0 {
			t.Errorf("handler testcase:\n%v\nnot update Time correctly. Initial value: %v current value: %v", test, initTime, htmlCounterCtx.Time)
		}

		// check the response body
		tpl, err := loadTemplate(*tplPath)
		if err != nil {
			t.Fatalf("error setting up response body check: %v", err)
		}
		buf := bytes.NewBuffer(make([]byte, 0))
		tpl.Execute(buf, htmlCounterCtx)

		if !bytes.Equal(buf.Bytes(), w.Body.Bytes()) {
			t.Errorf("handleGet: wrong body \nwant:\n\n%s\n\ngot:\n\n%s\n", buf.String(), w.Body.String())
		}
	}
}
