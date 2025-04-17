package app

import (
	"net"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(host, port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort(host, port)
	}
}
