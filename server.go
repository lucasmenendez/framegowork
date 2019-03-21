package shgf

import (
	"fmt"
	"net"
	"net/http"
)

const (
	minPort = 0
	maxPort = 65535
)

type server struct {
	address  string
	hostname string
	port     int
	debug    bool
	routes   routes
}

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

func (s *server) addRoute(r *route) error {
	if found := s.routes.exists(*r); found {
		return NewServerErr("route already registered")
	}

	s.routes = append(s.routes, r)
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.panic(w)

	var c *route
	if set := s.routes.findByMatchingPath(r.URL.Path); len(set) == 0 {
		NotFound().Submit(w)
		return
	} else if subset := set.findByMethod(r.Method); len(subset) == 0 {
		MethodNotAllowed().Submit(w)
		return
	} else {
		c = subset[0]
	}

	var ctx = &Context{route: c, Request: r}
	ctx.next = c.handler
	if err := c.handle(ctx).Submit(w); err != nil {
		InternalServerErr(err).Submit(w)
	}
}

func (s *server) start() error {
	if s.debug {
		fmt.Printf("Listen on: %s\n", s.address)
	}

	if e := http.ListenAndServe(s.address, s); e != nil {
		return NewServerErr("error starting server server", e)
	}
	return nil
}

func (s *server) panic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		InternalServerErr(r).Submit(w)
	}
}
