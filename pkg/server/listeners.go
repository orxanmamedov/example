package server

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type listeners struct {
	publicHTTP net.Listener
}

func (s *Server) newListener() (*listeners, error) {
	l := &listeners{}

	if s.cfg.HTTPPort != nil {
		publicHTTP, err := net.Listen("tcp", fmt.Sprintf(":%d", *s.cfg.HTTPPort))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create public HTTP port listener")
		}
		l.publicHTTP = publicHTTP
	}
	return l, nil
}
