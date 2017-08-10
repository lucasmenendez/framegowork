package framegowork

import (
	"fmt"
	"regexp"
	"strings"
	"net/http"
)

type Params map[string]string

//Router struct (class abstraction)
type Server struct {
	routes []Route
	port int
}

func New() *Server {
	return &Server{[]Route{}, 9999}
}

//Set port server
func (s *Server) SetPort(port int) {
	s.port = port
}

//GET
func (s *Server) GET(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	s.addMethod("GET", path, &handler, m)
}

//POST
func (s *Server) POST(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	s.addMethod("POST", path, &handler, m)
}

//PUT
func (s *Server) PUT(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	s.addMethod("PUT", path, &handler, m)
}

//DELETE
func (s *Server) DELETE(path string, handler Handler, middlewares ...Middleware) {
	var m *Middleware
	if len(middlewares) > 0 {
		m = &middlewares[0]
	}
	s.addMethod("DELETE", path, &handler, m)
}

//Create route with path, functions, methods, regexp to compare and middleware if exists
func (s *Server) addMethod(method, path string, handler *Handler, middleware *Middleware) {
	var position int = -1
	for i, route := range s.routes {
		if route.path == path {
			position = i
			break
		}
	}

	var res []string
	var pathComponents []string = strings.Split(path, "/")
	for _, str := range pathComponents {
		if len(str) > 0 && string(str[0]) == string(":") {
			str = "([A-Za-z0-9-_]+)"
		}
		res = append(res, str)
	}

	var rgx *regexp.Regexp = regexp.MustCompile(strings.Join(res, "/"))
	if position > -1 {
		s.routes[position].methods = append(s.routes[position].methods, method)
		s.routes[position].funcs = append(s.routes[position].funcs, handler)
		s.routes[position].rgx = rgx
		s.routes[position].middleware = middleware
	} else {
		var methods []string = []string{method}
		var funcs []*Handler = []*Handler{handler}
		s.routes = append(s.routes, Route{path, methods, funcs, rgx, middleware})
	}
}

//Listen routes and call, if exists, its function. Set router headers
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range s.routes {
		match, params := route.parsePath(r.URL.Path)
		if match {
			route.handleRoute(w, r, params)
			return
		}
	}
	http.Error(w, "Not found.", 404)
}

//Run on given port
func (s *Server) Run() {
	var port string = fmt.Sprintf(":%d", s.port)
	http.ListenAndServe(port, s)
}
