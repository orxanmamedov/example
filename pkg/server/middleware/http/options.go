package http

import (
	"net/http"

	"github.com/go-chi/cors"
)

type options struct {
	corsOptions cors.Options
	customMiddleware []func (http.Handler)  http.Handler
	opNameFunc operationNameFunc
}

type operationNameFunc func(*http.Request) string

func initOptions(opts []Option) *options {
	o := &options{}
	for i := range opts {
		opts[i](o)
	}
	return o
}

type Option func(*options)

func WithCORSOptions(corsOpts cors.Options) Option {
	return func (opts *options)  {
		opts.corsOptions = corsOpts
	}
}

func WithCustomMiddleware(mw ...func(http.Handler) http.Handler) Option {
	return func(opts *options) {
		opts.customMiddleware = mw
	}
}

func WithOperationNameFunc(f operationNameFunc) Option {
	return func (opts *options)  {
		opts.opNameFunc = f
	}
}