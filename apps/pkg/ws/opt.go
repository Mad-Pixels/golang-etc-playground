package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// Option wrapper for web server.
type Option func(ws *WS)

// WithMiddlewareReqID injects a request ID into the context (uid as string) of each request.
func WithMiddlewareReqID() Option {
	return func(ws *WS) {
		ws.router.Use(func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(),
					"uid",
					uuid.New().String(),
				)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
			return http.HandlerFunc(fn)
		})
	}
}

// WithMiddlewareRealIP add X-Real-IP header.
func WithMiddlewareRealIP() Option {
	return func(ws *WS) {
		ws.router.Use(middleware.RealIP)
	}
}

// WithMiddlewareRecover recovers from panics.
func WithMiddlewareRecover() Option {
	return func(ws *WS) {
		ws.router.Use(middleware.Recoverer)
	}
}

// WithMiddlewareDefaultLogger add default logger.
func WithMiddlewareDefaultLogger() Option {
	return func(ws *WS) {
		ws.router.Use(middleware.Logger)
	}
}

// WithMiddlewareCustomAfter add custom logger.
func WithMiddlewareCustomAfter(f func(w http.ResponseWriter, r *http.Request) error) Option {
	return func(ws *WS) {
		wrapper := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				if err := f(w, r); err != nil {
					panic(err)
				}
			})
		}
		ws.router.Use(wrapper)
	}
}

// WithMiddlewareCustomBefore add custom logger.
func WithMiddlewareCustomBefore(f func(w http.ResponseWriter, r *http.Request) error) Option {
	return func(ws *WS) {
		wrapper := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := f(w, r); err != nil {
					panic(err)
				}
				next.ServeHTTP(w, r)
			})
		}
		ws.router.Use(wrapper)
	}
}

// WithMiddlewareTimeout set request execution timeout.
func WithMiddlewareTimeout(sec int) Option {
	return func(ws *WS) {
		ws.router.Use(middleware.Timeout(time.Duration(1) * time.Second))
	}
}

// WithShutdownTimeout set timeout for shutting down listener.
func WithShutdownTimeout(sec int16) Option {
	return func(ws *WS) {
		ws.shutdownTimeout = sec
	}
}

// WithCustomPort set service port.
func WithCustomPort(p int) Option {
	return func(ws *WS) {
		ws.port = int32(p)
	}
}
