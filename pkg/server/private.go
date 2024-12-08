package server

import (
	"context"
	"time"

	"example/pkg/logger"
)

func (s *Server) closePrivate() error {
	// if s.monitoringHTTP == nil {
	// 	return errors.New("httpPrivate is nil")
	// }

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// s.monitoringHTTP.SetKeepAlivesEnabled(false)
	// if err := errors.Wrap(s.monitoringHTTP.Shutdown(ctx), "httpPrivate private server: error during shutdown"); err != nil {
	// 	return err
	// }
	logger.Warn(context.Background(), "httpPrivate private server: gracefully stopped")
	return nil
}
