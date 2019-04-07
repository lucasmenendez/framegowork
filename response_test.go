package shgf

import (
	"net/http/httptest"
	"testing"
)

func Test_parseBody(t *testing.T) {
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

func TestResponseJSON(t *testing.T) {
	var d = map[string]interface{}{
		"id":       1,
		"username": "testuser",
		"password": "testpass",
	}
	var rd = "{\"id\":1,\"password\":\"testpass\",\"username\":\"testuser\"}"

	var e error
	var r *Response
	if r, e = NewResponse(200); e != nil {
		t.Errorf("expected nil, got %s", e)
	}

	if e = r.JSON(d); e != nil {
		t.Errorf("expected nil, got %s", e)
	} else if h, ok := r.Header["Content-type"]; !ok {
		t.Error("expected true, got false")
	} else if h[0] != "application/json" {
		t.Errorf("expected \"application/json\", got %s", h[0])
	} else if b := string(r.Body); b != rd {
		t.Errorf("expected %s, got %s", rd, b)
	}
}

func TestResponse_submit(t *testing.T) {
	r, _ := NewResponse(200)
	w := httptest.NewRecorder()
	if e := r.submit(w, false); e != nil {
		t.Errorf("expected nil, got %s", e)
	}
}
