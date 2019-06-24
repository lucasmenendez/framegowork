package shgf

import (
	"net/http"
)

// Context struct contains the current request metainformation and utils like URL
// route params values, pointers to following functions or functions to handle
// with the request body.
type Context struct {
	route   *route
	next    *Handler
	Request *http.Request
	Params  Params
	Form    Form
}

// ParseParams function extract values of URL params defined by current route.
// Every param are labeled, so ParseParams only extract the values that match
// with route matcher naming each one with its label, and casting it with the
// type associated. All the params are stored into Params parameter of Context.
func (ctx *Context) ParseParams() (err error) {
	ctx.Params, err = decodeParams(ctx.Request.URL.Path, ctx.route.matcher)
	return
}

// ParseForm function invokes form function to parse the current request body to
// Form struct and append it into the current context. If an error occurs, it is
// returned.
func (ctx *Context) ParseForm() (err error) {
	ctx.Form, err = parseForm(ctx.Request)
	return
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
