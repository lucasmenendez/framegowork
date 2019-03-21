package shgf

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	pathRgx    = `(\/.*)+`
	paramDel   = `\/`
	paramTemp  = `<(int|string|float|bool):(.*)>`
	matcherRgx = `(?P<%s_%s>%s)`
	paramRgx   = `(int|string|float|bool)_(.+)`
	floatRgx   = `[0-9]+\.[0-9]+`
	intRgx     = `[0-9]+`
	boolRgx    = `true|false`
	strRgx     = `.+`
)

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

type Handler func(*Context) *Response

type route struct {
	method     string
	path       string
	middleware *Handler
	handler    *Handler
	matcher    *regexp.Regexp
}

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

func (r *route) parse() (e error) {
	if r.path == "" {
		e = NewServerErr("empty path not allowed")
		return
	} else if _, e = url.Parse(r.path); e != nil {
		return
	}

	var splitter, reader = regexp.MustCompile(paramDel), regexp.MustCompile(paramTemp)
	var items []string
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
