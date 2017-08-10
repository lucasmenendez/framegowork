package frameworkgo

import (
	"log"
	"regexp"
	"strings"
	"net/http"
)

//Route struct to store path with its methods and functions
type Route struct {
	path       string
	methods    []string
	funcs      []*Handler
	rgx        *regexp.Regexp
	middleware *Middleware
}

//Extract url params and check if route match with path
func (route Route) parsePath(path string) (bool, Params) {
	var attrs []string
	var params Params = make(Params)

	var routeComponents []string = strings.Split(route.path, "/")
	var pathComponents []string = strings.Split(path, "/")

	if len(routeComponents) == len(pathComponents) {
		for _, s := range routeComponents {
			if len(s) > 0 && string(s[0]) == string(":") {
				attrs = append(attrs, s[1:])
			}
		}

		if route.rgx.MatchString(path) {
			var values []string = route.rgx.FindStringSubmatch(path)[1:]
			for i, v := range values {
				params[attrs[i]] = v
			}
			return true, params
		} else {
			return false, params
		}
	} else {
		return false, params
	}
}

//Serve routes over all its methods
func (route Route) handleRoute(c Context) {
	for p, m := range route.methods {
		if m == c.Request.Method {
			if route.middleware == nil {
				f := *route.funcs[p]
				f(c)
				return
			} else {
				newContext := NewContext(c.Response, c.Request)
				newContext.Params = c.Params

				(*route.middleware)(newContext)
				return
			}
		}
	}
	http.Error(c.Response, "Not found.", 404)
}

func (route Route) handleRouteDebug(c Context) {
	log.Printf("[%s] %s", c.Request.Method, route.path)

	route.handleRoute(c)
}
