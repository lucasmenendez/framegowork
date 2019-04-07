package shgf

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	// minPort const contains the lower limit of valid port range.
	minPort = 0
	// maxPort const contains the upper limit of valid port range.
	maxPort = 65535
)

// server struct contains the current server configuration, attributes and
// registered routes into the same type.
type server struct {
	address  string
	hostname string
	port     int
	debug    bool
	routes   routes
}

// initServer function creates a new server by hostname and port provided. If
// debug mode is enabled, server will log errors. To create the new server,
// checks if the hostname provided and port are valids.
func initServer(h string, p int, d bool) (s *server, err error) {
	if ip := net.ParseIP(h); ip == nil {
		err = NewServerErr("invalid hostname IP")
		return
	} else if minPort >= p || p > maxPort {
		err = NewServerErr("port number out of bounds (0-65535)")
		return
	}

	return &server{
		address:  fmt.Sprintf("%s:%d", h, p),
		hostname: h,
		port:     p,
		debug:    d,
	}, nil
}

// addRoute function register route provided as new server function, ready to
// handle it. Before register it, checks if the route provided already exists
// into the current server.
func (s *server) addRoute(r *route) error {
	if found := s.routes.exists(*r); found {
		return NewServerErr("route already registered")
	}

	s.routes = append(s.routes, r)
	return nil
}

// ServerHTTP implements http.Handler interface to stay ready for
// http.ListAndServe function. The function catch all error into Internal Server
// Error. If debug is enabled, show traces of HTTP requests, then search into
// registered routes for matching one to call its handler, catching errors.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.panic(w)
	if s.debug {
		log.Printf("<- [%s] %s", r.Method, r.URL.Path)
	}

	var c *route
	if set := s.routes.findByMatchingPath(r.URL.Path); len(set) == 0 {
		NotFound().submit(w, s.debug)
		return
	} else if subset := set.findByMethod(r.Method); len(subset) == 0 {
		MethodNotAllowed().submit(w, s.debug)
		return
	} else {
		c = subset[0]
	}

	var ctx = &Context{route: c, Request: r}
	ctx.next = c.handler
	if err := c.handle(ctx).submit(w, s.debug); err != nil {
		InternalServerErr(err).submit(w, s.debug)
	}
}

// start function call to http.ListenAndServe function with current server
// instance.
func (s *server) start() error {
	if s.debug {
		log.Printf("Listen on: %s\n", s.address)
	}

	if e := http.ListenAndServe(s.address, s); e != nil {
		return NewServerErr("error starting server server", e)
	}
	return nil
}

// panic function is called if any runtime error occurs to recover it and
// wraps it into a Internal Server Error.
func (s *server) panic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		InternalServerErr(r).submit(w, s.debug)
	}
}
