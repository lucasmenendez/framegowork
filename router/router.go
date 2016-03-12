package router

import (
	"net/http"
	"regexp"
	"strings"
)

type Handler func(http.ResponseWriter, *http.Request, map[string]string)

//Route struct to store path with each methods and functions
type route struct {
	path    string
	methods []string
	funcs   []*Handler
}

//Router struct (class abstraction)
type Router struct {
	routes []route
}

//Contructor
func New() *Router {
	var routes []route
	return &Router{routes}
}

//Functions to add method and function to a path

//GET
func (r *Router) GET(path string, handler Handler) {
	r.addMethod("GET", path, &handler)
	return
}

//POST
func (r *Router) POST(path string, handler Handler) {
	r.addMethod("POST", path, &handler)
	return
}

//PUT
func (r *Router) PUT(path string, handler Handler) {
	r.addMethod("PUT", path, &handler)
	return
}

//DELETE
func (r *Router) DELETE(path string, handler Handler) {
	r.addMethod("DELETE", path, &handler)
	return
}

//Add Method check if path exists to append new method-function relation or create path with it
func (r *Router) addMethod(method, path string, handler *Handler) {
	position := -1

	for i, route := range r.routes {
		if route.path == path {
			position = i
			break
		}
	}

	if position > -1 {
		r.routes[position].methods = append(r.routes[position].methods, method)
		r.routes[position].funcs = append(r.routes[position].funcs, handler)
	} else {
		methods := []string{method}
		funcs := []*Handler{handler}
		r.routes = append(r.routes, route{path, methods, funcs})
	}

	return
}

//Extract url params and check if route match with path
func (route route) parsePath(path string) (bool, map[string]string) {
	var params, res []string

	attrs := make(map[string]string)

	route_strs := strings.Split(route.path, "/")
	path_strs := strings.Split(path, "/")

	if len(route_strs) == len(path_strs) {
		for _, str := range route_strs {
			if len(str) > 0 && string(str[0]) == string(":") {
				params = append(params, str[1:])
				str = "([A-Za-z0-9-_]+)"
			}
			res = append(res, str)
		}

		rgx, _ := regexp.Compile(strings.Join(res, "/"))
		if rgx.MatchString(path) {
			values := rgx.FindStringSubmatch(path)
			values = values[1:]
			for i, value := range values {
				attrs[params[i]] = value
			}

			return true, attrs
		} else {
			return false, attrs
		}
	} else {
		return false, attrs
	}
}

//Serve routes over all its methods
func (route route) handleRoute(w http.ResponseWriter, r *http.Request, params map[string]string) {
	for position, method := range route.methods {
		if method == r.Method {
			f := *route.funcs[position]
			f(w, r, params)
		}
	}
	return
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		match, params := route.parsePath(r.URL.Path)
		if match {
			route.handleRoute(w, r, params)
		}
	}
}

func (r *Router) Run(port string) {
	http.ListenAndServe(":"+port, r)
}
