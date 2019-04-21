package shgf

import (
	"net/http"
	"reflect"
	"testing"
)

type testRequest struct {
	valid    bool
	req      *http.Request
	expected *Form
}

func Test_parseForm(t *testing.T) {
	var validReq1, validReq2, invalidReq1, invalidReq2 *http.Request
	var form *Form

	var tests = []testRequest{
		{true, validReq1, form},
		{true, validReq2, form},
		{false, invalidReq1, nil},
		{false, invalidReq2, nil},
	}

	for _, test := range tests {
		if f, e := parseForm(test.req); e == nil && !test.valid {
			t.Error("expected error, got nil")
		} else if e != nil && test.valid {
			t.Errorf("expected nil, got %s", e)
		} else if !reflect.DeepEqual(test.expected, f) {
			t.Errorf("expected %+v, got %+v", test.expected, f)
		}
	}
}

func TestFormExists(t *testing.T) {}

func TestFormGet(t *testing.T) {}

func TestFormGetAll(t *testing.T) {}
