package shgf

import (
	"net/http"
	"regexp"
	"strconv"
)

// Context struct contains the current request metainformation and utils like URL
// route params values, pointers to following functions or functions to handle
// with the request body.
type Context struct {
	route   *route
	next    *Handler
	Request *http.Request
	Params  map[string]interface{}
	Form    *Form
}

// ParseParams function extract values of URL params defined by current route.
// Every param are labeled, so ParseParams only extract the values that match
// with route matcher naming each one with its label, and casting it with the
// type associated. All the params are stored into Params parameter of Context.
func (ctx *Context) ParseParams() (err error) {
	var (
		splitter = regexp.MustCompile(paramRgx)
		keys     = ctx.route.matcher.SubexpNames()
		values   = ctx.route.matcher.FindStringSubmatch(ctx.Request.URL.Path)
	)

	ctx.Params = map[string]interface{}{}
	keys, values = keys[1:], values[1:]
	for i, k := range keys {
		if k == "" || values[i] == "" {
			return NewServerErr("mismatch labeled params with values provided")
		}

		var metadata = splitter.FindStringSubmatch(k)[1:3]
		switch t, a, v := metadata[0], metadata[1], values[i]; t {
		case "float":
			ctx.Params[a], err = strconv.ParseFloat(v, 64)
			break
		case "int":
			ctx.Params[a], err = strconv.Atoi(v)
			break
		case "bool":
			ctx.Params[a] = (v == "true")
			break
		default:
			ctx.Params[a] = v
		}

		if err != nil {
			return
		}
	}

	return
}

// ParseForm function invokes form function to parse the current request body to
// Form struct and append it into the current context. If an error occurs, it is
// returned.
func (ctx *Context) ParseForm() (err error) {
	switch ctx.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodTrace, http.MethodPatch:
		ctx.Form, err = parseForm(ctx.Request)
		break
	default:
		err = NewServerErr("Current request method not allows to parse forms")
	}

	return err
}

// Next function invokes the main handler from middleware. If the Next function
// is invoked outside of middleware function, internal server error is returned.
func (ctx *Context) Next() *Response {
	if ctx.next == nil {
		var err = "next function not defined"
		return InternalServerErr(err)
	}

	var f = (*ctx.next)
	ctx.next = nil
	return f(ctx)
}
