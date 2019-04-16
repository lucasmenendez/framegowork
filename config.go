package shgf

import (
	"net"
	"os"
)

const (
	// minPort const contains the lower limit of valid port range.
	minPort = 0
	// maxPort const contains the upper limit of valid port range.
	maxPort = 65535
	// localPort const contains default HTTP port number for local services.
	localPort = 8080
	// localPortTLS const contains default HTTP port number for local services.
	localPortTLS = 8081
	// defaultPort const contains the localhost IP.
	localHostname = "127.0.0.1"
)

// Config struct allows to developers to configure the server easily. Requires
// hostname and port parameters and admits TLS cert and key paths and debug and
// HTTP2 flags to enable it.
type Config struct {
	Hostname          string
	Port, PortTLS     int
	Debug, HTTP2, TLS bool
	TLSCert, TLSKey   string
}

// LocalConf function returns basic configration for deploy server locally. Port
// number is an optional parameter. By default, 8080.
func LocalConf() *Config {
	var c = &Config{
		Hostname: localHostname,
		Port:     localPort,
		PortTLS:  localPortTLS,
		Debug:    true,
	}

	return c
}

// BasicConf function returns basic configration by hostname and port
// provided. Default configuration includes HTTP2 enabled and debug disabled.
func BasicConf(hostname string, port int) (*Config, error) {
	var c = &Config{
		Hostname: hostname,
		Port:     port,
	}

	return c, c.check()
}

// check function validates that configuration includes the required hostname
// and port parameters.
func (conf *Config) check() error {
	if ip := net.ParseIP(conf.Hostname); ip == nil {
		return NewServerErr("invalid hostname IP")
	}

	if minPort >= conf.Port || conf.Port > maxPort {
		return NewServerErr("port number out of bounds (0-65535)")
	}

	if conf.TLSCert != "" && conf.TLSKey != "" {
		var e error
		if _, e = os.Stat(conf.TLSCert); e != nil {
			return NewServerErr("error with TLSCert file path provided", e)
		}

		if _, e = os.Stat(conf.TLSKey); e != nil {
			return NewServerErr("error with TLSKey file path provided", e)
		}

		if minPort >= conf.PortTLS || conf.PortTLS > maxPort {
			return NewServerErr("TLS port number out of bounds (0-65535)")
		}

		if conf.PortTLS == conf.Port {
			return NewServerErr("TLS port and main port must be different")
		}

		conf.TLS = true
	}

	if conf.HTTP2 && !conf.TLS {
		return NewServerErr("HTTP2 requires TLS protocol enabled")
	}

	return nil
}
