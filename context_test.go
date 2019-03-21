package shgf

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"testing"
)

type testParseParams struct {
	ctx          *Context
	failedCtx    bool
	params       map[string]interface{}
	failedParams bool
}

type testNext struct {
	ctx    *Context
	failed bool
	res    *Response
}

func TestContext_parseParams(t *testing.T) {
	var errTests = []testParseParams{
		{
			ctx: &Context{
				route:   &route{matcher: regexp.MustCompile(`\/request\/(?P<int_id>[0-9]+)\/(e)dit`)},
				Request: &http.Request{URL: &url.URL{Path: "/request/12/edit"}},
			},
			failedCtx:    true,
			params:       map[string]interface{}{"id": 12},
			failedParams: false,
		},
		{
			ctx: &Context{
				route:   &route{matcher: regexp.MustCompile(`\/request\/(?P<int_id>[0-9]+)\/edit`)},
				Request: &http.Request{URL: &url.URL{Path: "/request/12/edit"}},
			},
			failedCtx:    false,
			params:       map[string]interface{}{"id": "12"},
			failedParams: true,
		},
		{
			ctx: &Context{
				route:   &route{matcher: regexp.MustCompile(`\/request\/(?P<int_id>[0-9]+)\/edit`)},
				Request: &http.Request{URL: &url.URL{Path: "/request/12/edit"}},
			},
			failedCtx:    false,
			params:       map[string]interface{}{"id": 12},
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
		} else if !test.failedParams && !equal {
			t.Errorf("expected %+v, got %+v", test.params, test.ctx.Params)
		}
	}
}

func TestContextAddNext(t *testing.T) {

}

func TestContextNext(t *testing.T) {
	resPass, _ := NewResponse(200)
	var nextPass Handler = func(ctx *Context) *Response {
		r, _ := NewResponse(200)
		return r
	}

	var tests = []testNext{
		{
			ctx:    &Context{},
			failed: true,
			res:    InternalServerErr("next function not defined"),
		},
		{
			ctx:    &Context{next: &nextPass},
			failed: false,
			res:    resPass,
		},
	}

	for _, test := range tests {
		res := test.ctx.Next()
		if test.failed && !reflect.DeepEqual(test.res, res) {
			t.Errorf("expected distinct to %+v, got %+v", test.res.Status, res.Status)
		} else if !test.failed && !reflect.DeepEqual(test.res.Status, res.Status) {
			t.Errorf("expected %+v, got %+v", test.res, res)
		}
	}
}
