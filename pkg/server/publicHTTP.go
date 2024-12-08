package server

import (
	"context"
	"example/pkg/logger"
	mwhttp "example/pkg/server/middleware/http"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
)

type RouterHTTP struct {
	Pattern string
	Handler chi.Router
}

func (s *Server) runHTTPPublic(r chi.Router) error {
	router := chi.NewMux()
	if s.cfg.Logging {
		router.Use(mwhttp.WithLogger(logger.Logger().Desugar()))
	}

	router.Use(
		mwhttp.Middleware(
			mwhttp.WithOperationNameFunc(nil),
			mwhttp.WithCORSOptions(cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{
					http.MethodHead,
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			}),
		)...,
	)

	if r != nil {
		router.Mount("/", r)
	}

	if len(router.Routes()) == 0 {
		router.HandleFunc("/", http.NotFound)
	}

	s.publicHTTP = &http.Server{
		Handler:           router,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTime,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	//Add HealthCheck

	return nil
}

func (s *Server) closeHTTPPublic() error {
	if s.publicHTTP == nil {
		return errors.New("httpPublic is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeOut)
	defer cancel()
	//add healthcheck
	logger.Warn(context.Background(), "httpPublic public server: waiting stop of traffic")
	time.Sleep(gracefulDelay)
	logger.Warn(context.Background(), "httpPublic public server: shutting down")

	s.publicHTTP.SetKeepAlivesEnabled(false)
	if err := errors.Wrap(s.publicHTTP.Shutdown(ctx), "httpPublic public server: error shutdown"); err != nil {
		return err
	}
	logger.Warn(context.Background(), "httpPublic public server: stopped")

	return nil
}
