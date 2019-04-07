// Package shgf its a simple http framework for go language. The framework
// provides simple API to create a HTTP server and a group of functions to
// register new available routes with its handler by HTTP method.
package shgf

import "net/http"

// Server struct contains the server's hostname and port. Server, also has a
// base server associated. Any Server instance has associated functions to
// register new routes with its handler bty HTTP method.
type Server struct {
	Hostname string
	Port     int
	base     *server
}

// New function creates a new Server instance with hostname and port provided by
// the user.
func New(hostname string, port int, debug bool) (srv *Server, err error) {
	var base *server
	if base, err = initServer(hostname, port, debug); err != nil {
		return
	}
	srv = &Server{hostname, port, base}
	return
}

// register function receives a HTTP method, route path and associated handlers,
// one at least. The function creates new route with that parameters (checking
// handler and middleware) and, if that route not currently exists, add it to
// the server.
func (srv *Server) register(method, path string, handlers ...Handler) error {
	if len(handlers) == 0 {
		return NewServerErr("not handler provided")
	}

	var (
		e       error
		r       *route
		handler = handlers[0]
	)
	if r, e = newRoute(method, path, &handler); e != nil {
		return e
	}

	if len(handlers) == 2 {
		r.middleware = &handlers[1]
	} else if len(handlers) > 2 {
		NewServerErr("only can add one middleware per route")
	}

	return srv.base.addRoute(r)
}

// Get function add new route on GET HTTP method with path and handlers
// provided.
func (srv *Server) Get(path string, handlers ...Handler) error {
	return srv.register(http.MethodGet, path, handlers...)
}

// Head function add new route on HEAD HTTP method with path and handlers
// provided.
func (srv *Server) Head(path string, handlers ...Handler) error {
	return srv.register(http.MethodHead, path, handlers...)
}

// Post function add new route on POST HTTP method with path and handlers
// provided.
func (srv *Server) Post(path string, handlers ...Handler) error {
	return srv.register(http.MethodPost, path, handlers...)
}

// Put function add new route on PUT HTTP method with path and handlers
// provided.
func (srv *Server) Put(path string, handlers ...Handler) error {
	return srv.register(http.MethodPut, path, handlers...)
}

// Patch function add new route on PATCH HTTP method with path and handlers
// provided.
func (srv *Server) Patch(path string, handlers ...Handler) error {
	return srv.register(http.MethodPatch, path, handlers...)
}

// Delete function add new route on DELETE HTTP method with path and handlers
// provided.
func (srv *Server) Delete(path string, handlers ...Handler) error {
	return srv.register(http.MethodDelete, path, handlers...)
}

// Connect function add new route on CONNECT HTTP method with path and handlers
// provided.
func (srv *Server) Connect(path string, handlers ...Handler) error {
	return srv.register(http.MethodConnect, path, handlers...)
}

// Options function add new route on OPTIONS HTTP method with path and handlers
// provided.
func (srv *Server) Options(path string, handlers ...Handler) error {
	return srv.register(http.MethodOptions, path, handlers...)
}

// Trace function add new route on TRACE HTTP method with path and handlers
// provided.
func (srv *Server) Trace(path string, handlers ...Handler) error {
	return srv.register(http.MethodTrace, path, handlers...)
}

// Listen function starts base server to listen new requests.
func (srv *Server) Listen() error {
	return srv.base.start()
}
