package server

import (
	"net/http"
	"regexp"
	"strings"
)

type Handler func(http.ResponseWriter, *http.Request, map[string]string)
type Middleware func(http.ResponseWriter, *http.Request, NextHandler)

//Server config
type Config struct {
	port    string
	headers map[string]string
}

//Route struct to store path with its methods and functions
type Route struct {
	path       string
	methods    []string
	funcs      []*Handler
	rgx        *regexp.Regexp
	middleware *Middleware
}

//Router struct (class abstraction)
type Server struct {
	routes []Route
	config Config
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

//Contructor
func New() *Server {
	return &Server{[]Route{}, Config{"9999", map[string]string{}}}
}

//Set port server
func (server *Server) SetPort(port string) {
	server.config.port = port
}

//Set server headers
func (server *Server) SetHeader(attr, value string) {
	server.config.headers[attr] = value
}

//Functions to add method and function to a path

//GET
func (r *Server) GET(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	r.addMethod("GET", path, &handler, m)
	return
}

//POST
func (r *Server) POST(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	r.addMethod("POST", path, &handler, m)
	return
}

//PUT
func (r *Server) PUT(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	r.addMethod("PUT", path, &handler, m)
	return
}

//DELETE
func (r *Server) DELETE(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	r.addMethod("DELETE", path, &handler, m)
	return
}

//Create route with path, functions, methods, regexp to compare and middleware if exists
func (server *Server) addMethod(method, path string, handler *Handler, middleware *Middleware) {
	position := -1

	for i, route := range server.routes {
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
		server.routes[position].methods = append(server.routes[position].methods, method)
		server.routes[position].funcs = append(server.routes[position].funcs, handler)
		server.routes[position].rgx = rgx
		server.routes[position].middleware = middleware
	} else {
		methods := []string{method}
		funcs := []*Handler{handler}
		server.routes = append(server.routes, Route{path, methods, funcs, rgx, middleware})
	}

	return
}

//Extract url params and check if route match with path
func (route Route) parsePath(path string) (bool, map[string]string) {
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
func (route Route) handleRoute(w http.ResponseWriter, r *http.Request, params map[string]string) {
	for position, method := range route.methods {
		if method == r.Method {
			if route.middleware == nil {
				f := *route.funcs[position]
				f(w, r, params)
				return
			} else {
				m := *route.middleware
				m(w, r, NextHandler{route.funcs[position], params})
				return
			}
		}
	}
	http.Error(w, "Not found.", 404)
	return
}

//Listen routes and call, if exists, its function. Set router headers
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for attr, value := range server.config.headers {
		w.Header().Set(attr, value)
	}

	for _, route := range server.routes {
		match, params := route.parsePath(r.URL.Path)
		if match {
			route.handleRoute(w, r, params)
			return
		}
	}
	http.Error(w, "Not found.", 404)
	return
}

//Run on given port
func (server *Server) Run() {
	http.ListenAndServe(":"+server.config.port, server)
}
