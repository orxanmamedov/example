package server

import (
	"context"
	"example/pkg/closer"
	"example/pkg/logger"
	"fmt"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	monitoringHTTPDefault uint = 8081
	readTimeout                = 60 * time.Second
	writeTime                  = 60 * time.Second
	readHeaderTimeout          = 60 * time.Second
)

type Server struct {
	cfg *Config

	publicHTTP *http.Server

	listeners *listeners

	publicCloser  *closer.Closer
	privateCloser *closer.Closer
}

type Config struct {
	Env  string
	Name string

	Swagger bool
	Logging bool

	HTTPPort *uint
}

func New(cfg *Config) (*Server, error) {
	s := &Server{
		cfg: cfg,

		publicCloser:  closer.New(syscall.SIGTERM, syscall.SIGINT),
		privateCloser: closer.New(),
	}

	l, err := s.newListener()
	if err != nil {
		return s, errors.WithMessage(err, "can't init listeners")
	}
	s.listeners = l

	return s, nil
}

const (
	gracefulDelay   = 5 * time.Second
	gracefulTimeOut = 10 * time.Second
)

func (s *Server) Run(routers []RouterHTTP) {
	logger.Errorf(context.Background(), "app started %s, env %s; %s", s.cfg.Name, s.cfg.Env, logPorts(s.cfg.HTTPPort, nil, 0))

	if s.cfg.HTTPPort != nil {
		go func(routers []RouterHTTP) {
			r := chi.NewRouter()
			if s.cfg.Swagger {
				r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))
			}
			r.Route("/api", func(r chi.Router) {
				r.Get("/status", routerCheck)
				for i := range routers {
					r.Mount(routers[i].Pattern, routers[i].Handler)
				}

			})

			if err := s.runHTTPPublic(r); err != nil {
				logger.Error(context.Background(), err.Error())
				s.publicCloser.CloseAll()
			}
		}(routers)
		s.publicCloser.Add(s.closeHTTPPublic)
	}

	s.publicCloser.Wait()
	closer.CloseAll()
	closer.Wait()
	s.privateCloser.CloseAll()

	_ = logger.Logger().Sync()
}

func routerCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("x-force", "with-you")
	ResponseJSON(w, r, nil)
}

func logPorts(httpPort, grpcPort *uint, monitoringPort uint) string {
	b := strings.Builder{}
	if httpPort != nil {
		b.WriteString(fmt.Sprintf("http port: %d; ", *httpPort))
	}
	if grpcPort != nil {
		b.WriteString(fmt.Sprintf("grpc port: %d; ", *grpcPort))
	}
	b.WriteString(fmt.Sprintf("monitoring port: %d; ", monitoringPort))
	return b.String()
}
