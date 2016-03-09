package router

import (
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

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

//Serve routes over all its methods
func handleRoute(path string, methods []string, funcs []*Handler) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		for position, method := range methods {
			if method == r.Method {
				f := *funcs[position]
				f(w, r)
			}
		}
	})
	return
}

//Iterate over routes and launch goroutine
func (r *Router) RunServer(port string) {
	for _, route := range r.routes {
		go handleRoute(route.path, route.methods, route.funcs)
	}

	http.ListenAndServe(":"+port, nil)
}
