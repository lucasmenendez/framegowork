package shgf

import (
	"reflect"
	"testing"
)

type testNewRoute struct {
	method, path string
	handler      *Handler
	valid        bool
	expected     *route
}

type testParseRoute struct {
	path, rgx string
	valid     bool
}

type testHandleRoute struct {
	handler, middleware *Handler
	valid               bool
}

var testHandler Handler = func(ctx *Context) *Response {
	res, _ := NewResponse(200)
	return res
}

var testMiddleware Handler = func(ctx *Context) *Response {
	return ctx.Next()
}

func TestNewRoute(t *testing.T) {
	var wrongHandler Handler
	var tests = []testNewRoute{
		{"GET", "/", &testHandler, true, &route{method: "GET", path: "/", handler: &testHandler}},
		{"HEAD", "/foo", &testHandler, true, &route{method: "HEAD", path: "/foo", handler: &testHandler}},
		{"POST", "/foo/bar", &testHandler, true, &route{method: "POST", path: "/foo/bar", handler: &testHandler}},
		{"PUT", "/foo/", &testHandler, true, &route{method: "PUT", path: "/foo/", handler: &testHandler}},
		{"PATCH", "/foo/bar/", &testHandler, true, &route{method: "PATCH", path: "/foo/bar/", handler: &testHandler}},
		{"DELETE", "/", &testHandler, true, &route{method: "DELETE", path: "/", handler: &testHandler}},
		{"CONNECT", "/", &testHandler, true, &route{method: "CONNECT", path: "/", handler: &testHandler}},
		{"OPTIONS", "/", &testHandler, true, &route{method: "OPTIONS", path: "/", handler: &testHandler}},
		{"TRACE", "/", &testHandler, true, &route{method: "TRACE", path: "/", handler: &testHandler}},
		{"wrong", "/", &testHandler, false, &route{}},
		{"GET", "", &testHandler, false, &route{}},
		{"GET", "/", &wrongHandler, false, &route{}},
	}

	for _, test := range tests {
		if r, e := newRoute(test.method, test.path, test.handler); e == nil {
			test.expected.matcher = r.matcher

			if test.valid && !reflect.DeepEqual(r, test.expected) {
				t.Errorf("expected %+v, got %+v", test.expected, r)
			} else if !test.valid {
				t.Errorf("expected %+v, got %+v", test.expected, r)
			}
		} else if test.valid {
			t.Errorf("expected nil, got %+v", e)
		}
	}
}

func TestRoute_parse(t *testing.T) {
	var tests = []testParseRoute{
		{"/item/child/12", `\/item\/child\/12`, true},
		{"/item/child/<int:id>", `\/item\/child\/(?P<int_id>[0-9]+)`, true},
		{"/item/child/<badformat>", ``, false},
		{"/item/child/<wrongtype:name>", ``, false},
	}

	var r = &route{}
	for _, test := range tests {
		r.path = test.path
		if e := r.parse(); test.valid && e != nil {
			t.Errorf("expected valid route, got error %s on path %s", e, r.path)
		} else if test.valid && test.rgx != r.matcher.String() {
			t.Errorf("expected %s matcher, got %s", test.rgx, r.matcher.String())
		}
	}
}

func TestRoute_handle(t *testing.T) {
	var tests = []testHandleRoute{
		{&testHandler, nil, true},
		{&testHandler, &testMiddleware, true},
		{nil, &testMiddleware, false},
		{nil, nil, false},
	}

	var ctx = &Context{}
	var route = &route{}
	res, _ := NewResponse(200)
	err := InternalServerErr("")
	for _, test := range tests {
		route.handler = test.handler
		route.middleware = test.middleware
		if r := route.handle(ctx); test.valid && r == nil {
			t.Errorf("expected request, got nil")
		} else if test.valid && (r.Status != res.Status || string(r.Body) != string(res.Body)) {
			t.Errorf("expected 200 response, got %d - %s", r.Status, string(r.Body))
		} else if !test.valid && (r != nil && r.Status != err.Status) {
			t.Errorf("expected %d - %s, got %d - %s", err.Status, err.Body, r.Status, r.Body)
		}
	}
}
