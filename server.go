package shgf

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

// server struct contains the current server configuration, attributes and
// registered routes into the same type.
type server struct {
	Addr, AddrTLS string
	conf          *Config
	routes        routes
}

// initServer function creates a new server by hostname and port provided. If
// debug mode is enabled, server will log errors. To create the new server,
// checks if the hostname provided and port are valids.
func initServer(c *Config) (s *server, err error) {
	return &server{
		Addr:    fmt.Sprintf("%s:%d", c.Hostname, c.Port),
		AddrTLS: fmt.Sprintf("%s:%d", c.Hostname, c.TLSPort),
		conf:    c,
	}, nil
}

// register function receives a HTTP method, route path and associated handlers,
// one at least. The function creates new route with that parameters (checking
// handler and middleware) and, if that route not currently exists, add it to
// the server. Before register it, checks if the route provided already exists
// into the current server.
func (s *server) register(method, path string, handlers ...Handler) error {
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
		return NewServerErr("only can add one middleware per route")
	}

	if found := s.routes.exists(*r); found {
		return NewServerErr("route already registered")
	}

	s.routes = append(s.routes, r)
	return nil
}

func (s *server) redirectToHTTPS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var uri = fmt.Sprintf("https://%s/%s", s.AddrTLS, r.RequestURI)
		http.Redirect(w, r, uri, http.StatusMovedPermanently)
	}
}

// ServerHTTP implements http.Handler interface to stay ready for
// http.ListAndServe function. The function catch all error into Internal Server
// Error. If debug is enabled, show traces of HTTP requests, then search into
// registered routes for matching one to call its handler, catching errors.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.panic(w)
	if s.conf.Debug {
		log.Printf("<- [%s] %s", r.Method, r.URL.Path)
	}

	var c *route
	if set := s.routes.findByMatchingPath(r.URL.Path); len(set) == 0 {
		NotFound().submit(w, s.conf.Debug)
		return
	} else if subset := set.findByMethod(r.Method); len(subset) == 0 {
		MethodNotAllowed().submit(w, s.conf.Debug)
		return
	} else {
		c = subset[0]
	}

	var ctx = &Context{route: c, Request: r}
	ctx.next = c.handler
	if err := c.handle(ctx).submit(w, s.conf.Debug); err != nil {
		InternalServerErr(err).submit(w, s.conf.Debug)
	}
}

// start function call to http.ListenAndServe function with current server
// instance.
func (s *server) start() error {
	if s.conf.Debug {
		if s.conf.TLS {
			log.Printf("Listen on: https://%s\n", s.AddrTLS)
		} else {
			log.Printf("Listen on: http://%s\n", s.Addr)
		}
	}

	if s.conf.TLS {
		if s.conf.HTTP2 {
			http2.ConfigureServer(&http.Server{Addr: s.Addr}, nil)
		}

		go func() {
			http.ListenAndServeTLS(s.AddrTLS, s.conf.TLSCert, s.conf.TLSKey, s)
		}()

		if e := http.ListenAndServe(s.Addr, s.redirectToHTTPS()); e != nil {
			return NewServerErr("error starting server")
		}

		return nil
	}

	if e := http.ListenAndServe(s.Addr, s); e != nil {
		return NewServerErr("error starting server", e)
	}

	return nil
}

// panic function is called if any runtime error occurs to recover it and
// wraps it into a Internal Server Error.
func (s *server) panic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		InternalServerErr(r).submit(w, s.conf.Debug)
	}
}
