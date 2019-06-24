package shgf

import (
	"reflect"
	"testing"
)

type testNext struct {
	ctx    *Context
	failed bool
	res    *Response
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
