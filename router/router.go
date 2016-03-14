package router

import (
	"net/http"
	"regexp"
	"strings"
)

type Handler func(http.ResponseWriter, *http.Request, map[string]string)
type Middleware func(http.ResponseWriter, *http.Request, NextHandler)

//Route struct to store path with its methods and functions
type route struct {
	path       string
	methods    []string
	funcs      []*Handler
	rgx        *regexp.Regexp
	middleware *Middleware
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

//Abstraction to exec next function
type NextHandler struct {
	handler *Handler
	params  map[string]string
}

//Exec next function when route has a middleware
func (n NextHandler) Exec(w http.ResponseWriter, r *http.Request) {
	next := *n.handler
	next(w, r, n.params)
	return
}

//Functions to add method and function to a path

//GET
func (r *Router) GET(path string, handler Handler, middlewares ...Middleware) {
	var m []*Middleware
	for _, middleware := range middlewares {
		m = append(m, &middleware)
	}
	r.addMethod("GET", path, &handler, m)
	return
}

//POST
func (r *Router) POST(path string, handler Handler, middlewares ...Middleware) {
	var m []*Middleware
	for _, middleware := range middlewares {
		m = append(m, &middleware)
	}
	r.addMethod("POST", path, &handler, m)
	return
}

//PUT
func (r *Router) PUT(path string, handler Handler, middlewares ...Middleware) {
	var m []*Middleware
	for _, middleware := range middlewares {
		m = append(m, &middleware)
	}
	r.addMethod("PUT", path, &handler, m)
	return
}

//DELETE
func (r *Router) DELETE(path string, handler Handler, middlewares ...Middleware) {
	var m []*Middleware
	for _, middleware := range middlewares {
		m = append(m, &middleware)
	}
	r.addMethod("DELETE", path, &handler, m)
	return
}

//Add Method check if path exists to append new method-function relation or create path with it
func (r *Router) addMethod(method, path string, handler *Handler, middlewares []*Middleware) {
	var middleware *Middleware
	if len(middlewares) > 0 {
		middleware = middlewares[0]
	}

	position := -1

	for i, route := range r.routes {
		if route.path == path {
			position = i
			break
		}
	}

	var res []string

	path_strs := strings.Split(path, "/")
	for _, str := range path_strs {
		if len(str) > 0 && string(str[0]) == string(":") {
			str = "([A-Za-z0-9-_]+)"
		}
		res = append(res, str)
	}

	rgx, _ := regexp.Compile(strings.Join(res, "/"))

	if position > -1 {
		r.routes[position].methods = append(r.routes[position].methods, method)
		r.routes[position].funcs = append(r.routes[position].funcs, handler)
		r.routes[position].rgx = rgx
		r.routes[position].middleware = middleware
	} else {
		methods := []string{method}
		funcs := []*Handler{handler}
		r.routes = append(r.routes, route{path, methods, funcs, rgx, middleware})
	}

	return
}

//Extract url params and check if route match with path
func (route route) parsePath(path string) (bool, map[string]string) {
	var params []string

	attrs := make(map[string]string)

	route_strs := strings.Split(route.path, "/")
	path_strs := strings.Split(path, "/")

	if len(route_strs) == len(path_strs) {

		for _, str := range route_strs {
			if len(str) > 0 && string(str[0]) == string(":") {
				params = append(params, str[1:])
			}
		}

		if route.rgx.MatchString(path) {
			values := route.rgx.FindStringSubmatch(path)
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
			if route.middleware == nil {
				f := *route.funcs[position]
				f(w, r, params)
			} else {
				m := *route.middleware
				m(w, r, NextHandler{route.funcs[position], params})
			}
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
