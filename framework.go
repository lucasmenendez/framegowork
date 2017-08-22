package frameworkgo

import (
	"fmt"
	"regexp"
	"net/http"
)

//Router struct (class abstraction)
type Server struct {
	routes	[]Route
	debug	bool
	port	int
}

func New() *Server {
	return &Server{[]Route{}, false, 9999}
}

//Set port server
func (s *Server) SetPort(port int) {
	s.port = port
}

//Set debug mode
func (s *Server) DebugMode(mode bool) {
	s.debug = mode
}

//GET
func (s *Server) GET(path string, handler Handler, middlewares ...Handler) {
	var mid *Handler
	if len(middlewares) > 0 {
		mid = &middlewares[0]
	}
	s.addMethod("GET", path, &handler, mid)
}

//POST
func (s *Server) POST(path string, handler Handler, middlewares ...Handler) {
	var mid *Handler
	if len(middlewares) > 0 {
		mid = &middlewares[0]
	}
	s.addMethod("POST", path, &handler, mid)
}

//PUT
func (s *Server) PUT(path string, handler Handler, middlewares ...Handler) {
	var mid *Handler
	if len(middlewares) > 0 {
		mid = &middlewares[0]
	}
	s.addMethod("PUT", path, &handler, mid)
}

//DELETE
func (s *Server) DELETE(path string, handler Handler, middlewares ...Handler) {
	var mid *Handler
	if len(middlewares) > 0 {
		mid = &middlewares[0]
	}
	s.addMethod("DELETE", path, &handler, mid)
}

//Create route with path, functions, methods, regexp to compare and middleware if exists
func (s *Server) addMethod(method, path string, handler *Handler, middleware *Handler) {
	var position int = -1
	for i, route := range s.routes {
		if route.path == path {
			position = i
			break
		}
	}

	var paramRegex *regexp.Regexp = regexp.MustCompile(`:([A-Za-z0-9-_]+)`)
	var newPath []byte = paramRegex.ReplaceAll([]byte(path), []byte("([A-Za-z0-9-_]+)"))
	var rgx *regexp.Regexp = regexp.MustCompile(string(newPath))

	if position > -1 {
		s.routes[position].methods = append(s.routes[position].methods, method)
		s.routes[position].handlers = append(s.routes[position].handlers, handler)
		s.routes[position].rgx = rgx
		s.routes[position].middleware = middleware
	} else {
		var methods []string = []string{method}
		var handlers []*Handler = []*Handler{handler}
		s.routes = append(s.routes, Route{path, methods, handlers, rgx, middleware})
	}
}

//Listen routes and call, if exists, its function. Set router headers
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range s.routes {
		match, params := ParseParams(route, r.URL.Path)
		if match {
			c := NewContext(r.URL.Path, w, r)
			c.Params = params
			if s.debug {
				route.handleRouteDebug(c)
			} else {
				route.handleRoute(c)
			}
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
