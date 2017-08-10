package frameworkgo

type Handler func(Context)
type Middleware func(Context)

//Abstraction to exec next function
type NextHandler struct {
	handler *Handler
	params  Params
}

//Exec next function when route has a middleware
func (n NextHandler) Exec(c Context) {
	var next Handler = *n.handler
	next(c)
}