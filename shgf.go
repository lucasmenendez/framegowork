// Package shgf its a simple http framework for go language. The framework
// provides simple API to create a HTTP server and a group of functions to
// register new available routes with its handler by HTTP method.
package shgf

import "net/http"

// Server struct contains the server's hostname and port. Server, also has a
// base server associated. Any Server instance has associated functions to
// register new routes with its handler by HTTP method.
type Server struct {
	base *server
}

// New function creates a new Server instance with config by the user. If debug
// mode is enabled, all inbound and outbound request will be logged.
func New(conf *Config) (srv *Server, err error) {
	if err = conf.check(); err != nil {
		return
	}

	var base *server
	if base, err = initServer(conf); err != nil {
		return
	}
	srv = &Server{base}
	return
}

// Get function add new route on GET HTTP method with path and handlers
// provided.
func (srv *Server) Get(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodGet, path, handlers...)
}

// Head function add new route on HEAD HTTP method with path and handlers
// provided.
func (srv *Server) Head(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodHead, path, handlers...)
}

// Post function add new route on POST HTTP method with path and handlers
// provided.
func (srv *Server) Post(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodPost, path, handlers...)
}

// Put function add new route on PUT HTTP method with path and handlers
// provided.
func (srv *Server) Put(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodPut, path, handlers...)
}

// Patch function add new route on PATCH HTTP method with path and handlers
// provided.
func (srv *Server) Patch(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodPatch, path, handlers...)
}

// Delete function add new route on DELETE HTTP method with path and handlers
// provided.
func (srv *Server) Delete(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodDelete, path, handlers...)
}

// Connect function add new route on CONNECT HTTP method with path and handlers
// provided.
func (srv *Server) Connect(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodConnect, path, handlers...)
}

// Options function add new route on OPTIONS HTTP method with path and handlers
// provided.
func (srv *Server) Options(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodOptions, path, handlers...)
}

// Trace function add new route on TRACE HTTP method with path and handlers
// provided.
func (srv *Server) Trace(path string, handlers ...Handler) error {
	return srv.base.register(http.MethodTrace, path, handlers...)
}

// Listen function starts base server to listen new requests.
func (srv *Server) Listen() error {
	return srv.base.start()
}
