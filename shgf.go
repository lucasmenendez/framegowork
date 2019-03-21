package shgf

import (
	"net/http"
)

type Server struct {
	Hostname string
	Port     int
	Debug    bool
	base     *server
}

func New(hostname string, port int, debug bool) (srv *Server, err error) {
	var base *server
	if base, err = initServer(hostname, port, debug); err != nil {
		return
	}
	srv = &Server{hostname, port, debug, base}
	return
}

func (srv *Server) register(method, path string, handlers ...Handler) error {
	if len(handlers) == 0 {
		return NewServerErr("not handler provided")
	}

	var handler = handlers[0]
	if r, e := newRoute(method, path, &handler); e != nil {
		return e
	} else {
		if len(handlers) == 2 {
			r.middleware = &handlers[1]
		} else if len(handlers) > 2 {
			NewServerErr("only can add one middleware per route")
		}

		return srv.base.addRoute(r)
	}
}

func (srv *Server) GET(path string, handlers ...Handler) error {
	return srv.register(http.MethodGet, path, handlers...)
}

func (srv *Server) HEAD(path string, handlers ...Handler) error {
	return srv.register(http.MethodHead, path, handlers...)
}

func (srv *Server) POST(path string, handlers ...Handler) error {
	return srv.register(http.MethodPost, path, handlers...)
}

func (srv *Server) PUT(path string, handlers ...Handler) error {
	return srv.register(http.MethodPut, path, handlers...)
}

func (srv *Server) PATCH(path string, handlers ...Handler) error {
	return srv.register(http.MethodPatch, path, handlers...)
}

func (srv *Server) DELETE(path string, handlers ...Handler) error {
	return srv.register(http.MethodDelete, path, handlers...)
}

func (srv *Server) CONNECT(path string, handlers ...Handler) error {
	return srv.register(http.MethodConnect, path, handlers...)
}

func (srv *Server) OPTIONS(path string, handlers ...Handler) error {
	return srv.register(http.MethodOptions, path, handlers...)
}

func (srv *Server) TRACE(path string, handlers ...Handler) error {
	return srv.register(http.MethodTrace, path, handlers...)
}

func (srv *Server) Listen() error {
	return srv.base.start()
}
