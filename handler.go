package frameworkgo

type Handler func(Response, Request, Params)
type Middleware func(Response, Request, NextHandler)

//Abstraction to exec next function
type NextHandler struct {
	handler *Handler
	params  Params
}

//Exec next function when route has a middleware
func (n NextHandler) Exec(w Response, r Request) {
	var next Handler = *n.handler
	next(w, r, n.params)
}

