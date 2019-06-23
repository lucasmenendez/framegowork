package shgf

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type testRequest struct {
	valid    bool
	req      *http.Request
	expected *Form
}

func Test_parseForm(t *testing.T) {
	var form = &Form{
		keys:   []string{"foo"},
		fields: map[string]interface{}{"foo": []string{"bar"}},
	}

	var params = url.Values{}
	params.Set("foo", "bar")
	var data = params.Encode()
	var validReq1 = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	validReq1.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var buf = bytes.NewBuffer([]byte{})
	var writer = multipart.NewWriter(buf)
	if label, e := writer.CreateFormField("foo"); e == nil {
		label.Write([]byte("bar"))
		writer.Close()
	} else {
		t.Errorf("expected nil, got %s", e)
		return
	}
	var validReq2 = httptest.NewRequest(http.MethodPost, "/", buf)
	validReq2.Header.Set("Content-Type", writer.FormDataContentType())

	var invalidReq1 = httptest.NewRequest(http.MethodGet, "/", nil)
	invalidReq1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var invalidReq2 = httptest.NewRequest(http.MethodPut, "/", nil)
	invalidReq2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var tests = []testRequest{
		{true, validReq1, form},
		{true, validReq2, form},
		{false, invalidReq1, form},
		{false, invalidReq2, form},
	}

	for _, test := range tests {
		if f, e := parseForm(test.req); test.valid {
			if e != nil {
				t.Errorf("expected nil, got %s", e)
			}

			if !reflect.DeepEqual(f, form) {
				t.Errorf("expected %+v, got %+v", form, f)
			}
		} else if e == nil && reflect.DeepEqual(f, form) {
			t.Error("expected error, got nil")
		}
	}
}

func TestFormExists(t *testing.T) {
	var form = &Form{
		keys:   []string{"foo"},
		fields: map[string]interface{}{"foo": []string{"bar"}},
	}

	if exists := form.Exists("foo"); !exists {
		t.Error("expected true, got false")
	}

	if exists := form.Exists("invalid"); exists {
		t.Error("expected false, got true")
	}

	if exists := form.Exists(""); exists {
		t.Error("expected false, got true")
	}
}

func TestFormGet(t *testing.T) {
	var form = &Form{
		keys:   []string{"foo"},
		fields: map[string]interface{}{"foo": []string{"bar"}},
	}

	if res, err := form.Get("foo"); err != nil {
		t.Errorf("expected nil, got %s", err)
	} else if !reflect.DeepEqual(res, []string{"bar"}) {
		t.Errorf("expected %+v, got %+v", []string{"bar"}, res)
	}

	if _, err := form.Get("invalid"); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFormGetAll(t *testing.T) {
	var form = &Form{
		keys:   []string{"foo"},
		fields: map[string]interface{}{"foo": []string{"bar"}},
	}

	if res := form.GetAll(); !reflect.DeepEqual(res, map[string]interface{}{"foo": []string{"bar"}}) {
		t.Errorf("expected %+v, got %+v", map[string]interface{}{"foo": []string{"bar"}}, res)
	}

	form = &Form{
		keys:   []string{},
		fields: map[string]interface{}{},
	}

	if res := form.GetAll(); !reflect.DeepEqual(res, map[string]interface{}{}) {
		t.Errorf("expected %+v, got %+v", map[string]interface{}{}, res)
	}
}
