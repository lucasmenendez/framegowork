package shgf

const (
	minPort = 0
	maxPort = 65535
)

type server struct {
	hostname string
	port     int
	debug    bool
	routes   routes
}

func newServer(h string, p int, d ...bool) (s *server, e error) { return s, e }
func (s *server) isReady() bool                                 { return true }
func (s *server) addRoute(m, p string, h handler) error         { return nil }
func (s *server) start() error                                  { return nil }
func (s *server) stop() error                                   { return nil }
