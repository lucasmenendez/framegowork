package router

import (
	"net/http"
	"reflect"
)

//Basic http function
type Handle func(http.ResponseWriter, *http.Request)

//Route struct to store path with each methods and functions
type route struct {
	path    string
	methods []string
	funcs   []Handle
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
func (r *Router) GET(path string, handle Handle) {
	r.addMethod("GET", path, handle)
	return
}

//POST
func (r *Router) POST(path string, handle Handle) {
	r.addMethod("POST", path, handle)
	return
}

//PUT
func (r *Router) PUT(path string, handle Handle) {
	r.addMethod("PUT", path, handle)
	return
}

//DELETE
func (r *Router) DELETE(path string, handle Handle) {
	r.addMethod("DELETE", path, handle)
	return
}


//Add Method check if path exists to append new method-function relation or create path with it
func (r *Router) addMethod(method, path string, handle Handle) {
	position := -1

	for i, route := range r.routes {
		if route.path == path {
			position = i
			break
		}
	}

	if position > -1 {
		r.routes[position].methods = append(r.routes[position].methods, method)
		r.routes[position].funcs = append(r.routes[position].funcs, handle)
	} else {
		methods := []string{method}
		funcs := []Handle{handle}
		r.routes = append(r.routes, route{path, methods, funcs})
	}

	return
}

//Serve routes over all its methods
func handleRoute(path string, methods []string, funcs []Handle) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		for position, method := range methods {
			if method == r.Method {
				params := make([]reflect.Value, 2)
				params[0] = reflect.ValueOf(w)
				params[1] = reflect.ValueOf(r)

				f := reflect.ValueOf(funcs[position])
				f.Call(params)
			}
		}
	})
}

//Iterate over routes and launch goroutine
func (r *Router) RunServer(port string) {
	for _, route := range r.routes {
		go handleRoute(route.path, route.methods, route.funcs)
	}
	http.ListenAndServe(":"+port, nil)
}
