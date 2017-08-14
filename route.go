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
		if m == c.request.Method {
			if route.middleware == nil {
				f := *route.handlers[p]
				f(c)
				return
			} else {
				newContext := NewContext(route.path, c.response, c.request)
				newContext.Params = c.Params
				newContext.Handler = *route.handlers[p]

				(*route.middleware)(newContext)
				return
			}
		}
	}
	http.Error(c.response, "Method not allowed.", 405)
}

func (route Route) handleRouteDebug(c Context) {
	log.Printf("[%s] %s", c.request.Method, c.Path)
	route.handleRoute(c)
}
