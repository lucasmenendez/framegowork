package shgf

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"testing"
)

type testParamsEncoded struct {
	path, rgx string
	valid     bool
}

type testParamsDecoded struct {
	ctx          *Context
	failedCtx    bool
	params       Params
	failedParams bool
}

func TestParams_encodeParams(t *testing.T) {
	var tests = []testParamsEncoded{
		{"/item/child/12", `\/item\/child\/12`, true},
		{"/item/child/<int:id>", `\/item\/child\/(?P<int_id>[0-9]+)`, true},
		{"/item/child/<badformat>", ``, false},
		{"/item/child/<wrongtype:name>", ``, false},
	}

	var r = &route{}
	for _, test := range tests {
		if rgx, e := encodeParams(test.path); test.valid && e != nil {
			t.Errorf("expected valid route, got error %s on path %s", e, r.path)
		} else if test.valid && rgx.String() != test.rgx {
			t.Errorf("expected %s matcher, got %s", test.rgx, rgx.String())
		}
	}
}

func TestParams_decodeParams(t *testing.T) {
	var errTests = []testParamsDecoded{
		{
			ctx: &Context{
				route:   &route{matcher: regexp.MustCompile(`\/request\/(?P<int_id>[0-9]+)\/(e)dit`)},
				Request: &http.Request{URL: &url.URL{Path: "/request/12/edit"}},
			},
			failedCtx:    true,
			params:       Params{"id": 12},
			failedParams: false,
		},
		{
			ctx: &Context{
				route:   &route{matcher: regexp.MustCompile(`\/request\/(?P<int_id>[0-9]+)\/edit`)},
				Request: &http.Request{URL: &url.URL{Path: "/request/12/edit"}},
			},
			failedCtx:    false,
			params:       Params{"id": "12"},
			failedParams: true,
		},
		{
			ctx: &Context{
				route:   &route{matcher: regexp.MustCompile(`\/request\/(?P<int_id>[0-9]+)\/edit`)},
				Request: &http.Request{URL: &url.URL{Path: "/request/12/edit"}},
			},
			failedCtx:    false,
			params:       Params{"id": 12},
			failedParams: false,
		},
	}

	for _, test := range errTests {
		if err := test.ctx.ParseParams(); test.failedCtx && err == nil {
			t.Errorf("error expected, got nil")
		} else if !test.failedCtx && err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		} else if equal := reflect.DeepEqual(test.params, test.ctx.Params); test.failedParams && equal {
			t.Errorf("expected distinct to %+v, got %+v", test.params, test.ctx.Params)
		} else if !test.failedParams && (!equal && err == nil) {
			t.Errorf("expected %+v, got %+v", test.params, test.ctx.Params)
		}
	}
}

func TestParamsExists(t *testing.T) {
	var params = Params{"foo": []string{"bar"}}

	if exists := params.Exists("foo"); !exists {
		t.Error("expected true, got false")
	}

	if exists := params.Exists("invalid"); exists {
		t.Error("expected false, got true")
	}

	if exists := params.Exists(""); exists {
		t.Error("expected false, got true")
	}
}

func TestParamsGet(t *testing.T) {
	var params = Params{"foo": []string{"bar"}}

	if res, err := params.Get("foo"); err != nil {
		t.Errorf("expected nil, got %s", err)
	} else if !reflect.DeepEqual(res, []string{"bar"}) {
		t.Errorf("expected %+v, got %+v", []string{"bar"}, res)
	}

	if _, err := params.Get("invalid"); err == nil {
		t.Error("expected error, got nil")
	}
}
