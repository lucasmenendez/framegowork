package frameworkgo

import (
	"log"
	"regexp"
	"net/http"
)

//Route struct to store path with its methods and functions
type Route struct {
	path       string
	methods    []string
	handlers   []*Handler
	rgx        *regexp.Regexp
	middleware *Handler
}

//Serve routes over all its methods
func (route Route) handleRoute(c Context) {
	for p, m := range route.methods {
		if m == c.Request.Method {
			if route.middleware == nil {
				f := *route.handlers[p]
				f(c)
				return
			} else {
				newContext := NewContext(route.path, c.Response, c.Request)
				newContext.Params = c.Params
				newContext.Handler = *route.handlers[p]

				(*route.middleware)(newContext)
				return
			}
		}
	}
	http.Error(c.Response, "Method not allowed.", 405)
}

func (route Route) handleRouteDebug(c Context) {
	log.Printf("[%s] %s", c.Request.Method, c.Path)
	route.handleRoute(c)
}
