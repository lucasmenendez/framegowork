package shgf

import (
	"net/http"
	"net/url"
	"regexp"
)

// validMethods contains all HTTP Methods into a string slice.
var validMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// Handler type wraps handler function that receives a Context pointer and
// returns a Response pointer as a result of handled Request.
type Handler func(*Context) *Response

// route struct contains required data to register a HTTP route to handle.
type route struct {
	method     string
	path       string
	middleware *Handler
	handler    *Handler
	matcher    *regexp.Regexp
}

// newRoute function create a new route by method, path and route handler
// provided. Checks if params provided are valid and call some other functions
// to prepare the route to be appended.
func newRoute(method, path string, handler *Handler) (r *route, e error) {
	var valid = false
	for _, validMethod := range validMethods {
		if valid = method == validMethod; valid {
			break
		}
	}

	if !valid {
		e = NewServerErr("HTTP method not allowed")
		return
	}

	if _, e = url.Parse(path); e != nil {
		e = NewServerErr("wrong path provided", e)
		return
	} else if path == "" {
		e = NewServerErr("empty path provided", e)
	}

	if *handler == nil {
		e = NewServerErr("nil handler not allowed")
		return
	}

	r = &route{
		method:  method,
		path:    path,
		handler: handler,
	}
	if e = r.parse(); e != nil {
		e = NewServerErr("error parsing route", e)
	}

	return
}

// parse function generates a matcher Regexp for the current route, including
// route params. Validates that the current route has the required data
// available and parses route params to generate the matcher. If any of this
// processes raise an error, function return ServerErr.
func (r *route) parse() error {
	if r.path == "" {
		return NewServerErr("empty path not allowed")
	} else if _, e := url.Parse(r.path); e != nil {
		return NewServerErr("wrong or bad formatted path provided", e)
	}

	var err error
	r.matcher, err = encodeParams(r.path)
	return err
}

// handle function calls to the current route handler or middleware, if it is
// available, with the context provided as a param. If any of both functions is
// available, returns nil, else returns the result of calling the function.
func (r *route) handle(ctx *Context) *Response {
	var f Handler
	if r.middleware != nil && r.handler != nil {
		f = *r.middleware
		ctx.next = r.handler
	} else if r.handler != nil {
		f = *r.handler
	}

	if f != nil {
		return f(ctx)
	}

	return nil
}
