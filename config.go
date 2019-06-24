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
	// localTLSPort const contains default HTTP port number for local services.
	localTLSPort = 8081
	// defaultPort const contains the localhost IP.
	localHostname = "127.0.0.1"
)

// Config struct allows to developers to configure the server easily. Requires
// hostname and port parameters and admits TLS cert and key paths and debug and
// HTTP2 flags to enable it.
type Config struct {
	Hostname          string
	Port, TLSPort     int
	Debug, HTTP2, TLS bool
	TLSCert, TLSKey   string
}

// LocalConf function returns basic configration for deploy server locally. Port
// number is an optional parameter. By default, 8080.
func LocalConf() *Config {
	var c = &Config{
		Hostname: localHostname,
		Port:     localPort,
		TLSPort:  localTLSPort,
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

	return conf.validTLS()
}

// validTLS function validates that configuration includes the required TLS
// parameters.
func (conf *Config) validTLS() error {
	if conf.TLSCert != "" && conf.TLSKey != "" {
		if minPort >= conf.TLSPort || conf.TLSPort > maxPort {
			return NewServerErr("TLS port number out of bounds (0-65535)")
		} else if conf.TLS && conf.TLSPort == conf.Port {
			return NewServerErr("TLS port and main port must be different")
		}

		if _, e := os.Stat(conf.TLSCert); e != nil {
			return NewServerErr("error with TLSCert file path provided", e)
		}

		if _, e := os.Stat(conf.TLSKey); e != nil {
			return NewServerErr("error with TLSKey file path provided", e)
		}

		conf.TLS = true
	} else if conf.HTTP2 {
		return NewServerErr("HTTP2 requires TLS configuration")
	}

	return nil
}
