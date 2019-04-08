package shgf

import "net"

const (
	// minPort const contains the lower limit of valid port range.
	minPort = 0
	// maxPort const contains the upper limit of valid port range.
	maxPort = 65535
)

// Config struct allows to developers to configure the server easily. Requires
// hostname and port parameters and admits TLS cert and key paths and debug and
// HTTP2 flags to enable it.
type Config struct {
	Hostname        string
	Port            int
	Debug, HTTP2    bool
	TLSCert, TLSKey string
}

// DefaultConf function returns basic configration by hostname and port
// provided. Default configuration includes HTTP2 enabled and debug disabled.
func DefaultConf(hostname string, port int) (Config, error) {
	var c = Config{
		Hostname: hostname,
		Port:     port,
		Debug:    false,
		HTTP2:    true,
	}

	return c, c.check()
}

// check function validates that configuration includes the required hostname
// and port parameters.
func (conf Config) check() (err error) {
	if ip := net.ParseIP(conf.Hostname); ip == nil {
		err = NewServerErr("invalid hostname IP")
		return
	} else if minPort >= conf.Port || conf.Port > maxPort {
		err = NewServerErr("port number out of bounds (0-65535)")
	}

	return
}
