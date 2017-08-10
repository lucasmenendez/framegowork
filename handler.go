package frameworkgo

import "net/http"

type Handler func(http.ResponseWriter, *http.Request, Params)
type Middleware func(http.ResponseWriter, *http.Request, NextHandler)

//Abstraction to exec next function
type NextHandler struct {
	handler *Handler
	params  Params
}

//Exec next function when route has a middleware
func (n NextHandler) Exec(w http.ResponseWriter, r *http.Request) {
	var next Handler = *n.handler
	next(w, r, n.params)
}

