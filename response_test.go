package shgf

import (
	"net/http"
	"testing"
)

type testResponse struct {
	t      *testing.T
	header http.Header
	body   []byte
	status int
}

func NewTestResponse(t *testing.T) *testResponse {
	return &testResponse{
		t:      t,
		header: make(http.Header),
	}
}

func (r *testResponse) Header() http.Header {
	return r.header
}

func (r *testResponse) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *testResponse) WriteHeader(status int) {
	r.status = status
}

func (r *testResponse) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != body {
		r.t.Errorf("expected body %+v to equal %+v", string(r.body), body)
	}
}

func TestParseBody(t *testing.T) {
	if _, e := parseBody(12); e == nil {
		t.Error("expected error, got nil")
	}

	if _, e := parseBody(1.3); e == nil {
		t.Error("expected error, got nil")
	}

	if _, e := parseBody(true); e == nil {
		t.Error("expected error, got nil")
	}

	if _, e := parseBody(map[string]interface{}{"id": 2, "name": "foo"}); e == nil {
		t.Error("expected error, got nil")
	}

	if _, e := parseBody([]byte("{\"id\":2,\"name\":\"foo\"}")); e != nil {
		t.Errorf("expected nil, got %s", e)
	}

	if _, e := parseBody("ok"); e != nil {
		t.Errorf("expected nil, got %s", e)
	}
}
func TestNewResponse(t *testing.T) {
	if r, e := NewResponse(534); e == nil {
		t.Error("expected error, got nil")
	} else if r.Status != 500 {
		t.Errorf("expected error 500, got %d", r.Status)
	}

	if _, e := NewResponse(200); e != nil {
		t.Errorf("expected nil, got %s", e)
	}

	if r, e := NewResponse(200, "All right", "Cause error"); e == nil {
		t.Error("expected error, got nil")
	} else if r.Status != 500 {
		t.Errorf("expected error 500, got %d", r.Status)
	}
}
func testResponse_submit(t *testing.T) {
	r, _ := NewResponse(200)
	if e := r.Submit(NewTestResponse(t)); e != nil {
		t.Errorf("expected nil, got %s", e)
	}
}
