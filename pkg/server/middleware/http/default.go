package http

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func Middleware(opts ...Option) []func(http.Handler) http.Handler {
	o := initOptions(opts)

	c := cors.New(o.corsOptions)

	mw := []func(http.Handler) http.Handler{
		defaultMiddleware(o.opNameFunc),
		mwGzipRequest,
		c.Handler,
		Recover,
		middleware.URLFormat,
		render.SetContentType(render.ContentTypeJSON),
	}

	mw = append(mw, o.customMiddleware...)
	return mw
}

func WithLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info(fmt.Sprintf("finished http request with code: %d", ww.Status()),
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
				)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func defaultMiddleware(opNameFunc operationNameFunc) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return newWrapper(h, opNameFunc)
	}
}

type defaultsWrapper struct {
	next http.Handler

	pool *sync.Pool

	opNameFunc operationNameFunc
}

func newWrapper(h http.Handler, opNameFunc operationNameFunc) *defaultsWrapper {
	w := &defaultsWrapper{
		next:       h,
		pool:       &sync.Pool{},
		opNameFunc: opNameFunc,
	}
	w.pool.New = func() interface{} {
		return chi.NewRouteContext()
	}
	if w.opNameFunc == nil {
		w.opNameFunc = w.getOpName
	}

	return w
}

func (h *defaultsWrapper) getOpName(r *http.Request) string {
	var opName, handler string
	if r != nil {
		opName, handler = "HTTP "+r.Method, ""
	}

	if r != nil && r.Method != "OPTIONS" && r.URL != nil {
		handler = "undefined"

		rctx := chi.RouteContext(r.Context())

		pctx := h.pool.Get().(*chi.Context)
		pctx.Reset()

		path := rctx.RoutePath
		if path == "" {
			if r.URL.RawPath != "" {
				path = r.URL.RawPath
			} else {
				path = r.URL.Path
			}
		}

		hasTrailingSlash := len(path) > 1 && path[len(path)-1] == '/'
		if hasTrailingSlash {
			path = path[:len(path)-1]
		}

		if ok := rctx.Routes.Match(pctx, r.Method, path); ok {
			handler = pctx.RoutePattern()
			if hasTrailingSlash {
				handler += "/"
			}
		}
		h.pool.Put(pctx)
	}

	if handler != "" {
		opName += ": " + handler
	}

	return opName
}

func (h *defaultsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opName := h.opNameFunc(r)
	if opName == "" {
		opName = h.getOpName(r)
	}

	h.next.ServeHTTP(w, r)
}

func mwGzipRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			var err error
			r.Body, err = gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
		}
		next.ServeHTTP(w, r)
	})
}
