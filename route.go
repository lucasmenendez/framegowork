package shgf

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	// pathRgx matches with any string between slashes. Used for validation.
	pathRgx = `(\/.*)+`
	// paramDel matches with slashes. Used for split and join paths.
	paramDel = `\/`
	// paramTemp matches with param format. Used for extract params from the
	// path.
	paramTemp = `<(int|string|float|bool):(.*)>`
	// matcherRgx its a template to generate the matcher to handling routes.
	matcherRgx = `(?P<%s_%s>%s)`
	// paramRgx mathes with the route handled param extracting param type and
	// content at once.
	paramRgx = `(int|string|float|bool)_(.+)`
	// floatRgx matches with float number into a string. Used to generate
	// the matcher to handling routes with float params.
	floatRgx = `[0-9]+\.[0-9]+`
	// intRgx matches with int number into a string. Used to generate
	// the matcher to handling routes with int params.
	intRgx = `[0-9]+`
	// boolRgx matches with bool number into a string. Used to generate
	// the matcher to handling routes with bool params.
	boolRgx = `true|false`
	// stringRgx matches with string number into a string. Used to generate
	// the matcher to handling routes with string params.
	strRgx = `.+`
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
func (r *route) parse() (e error) {
	if r.path == "" {
		e = NewServerErr("empty path not allowed")
		return
	} else if _, e = url.Parse(r.path); e != nil {
		return
	}

	var (
		splitter = regexp.MustCompile(paramDel)
		reader   = regexp.MustCompile(paramTemp)
		items    []string
	)
	for _, i := range splitter.Split(r.path, -1) {
		var item string
		if res := reader.FindStringSubmatch(i); len(res) == 3 {
			switch t, a := res[1], res[2]; t {
			case "float":
				item = fmt.Sprintf(matcherRgx, t, a, floatRgx)
				break
			case "int":
				item = fmt.Sprintf(matcherRgx, t, a, intRgx)
				break
			case "bool":
				item = fmt.Sprintf(matcherRgx, t, a, boolRgx)
				break
			case "string":
				item = fmt.Sprintf(matcherRgx, t, a, strRgx)
				break
			default:
				e = NewServerErr("unknown data type")
				return
			}
		} else {
			item = i
		}

		items = append(items, item)
	}

	var matcherTemp = strings.Join(items, paramDel)
	r.matcher = regexp.MustCompile(matcherTemp)
	return
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
