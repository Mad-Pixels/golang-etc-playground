package ws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apps "github.com/Mad-Pixels/golang-playground"

	"github.com/go-chi/chi/v5"
	"github.com/valyala/fasthttp/reuseport"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

// WS object.
type WS struct {
	router          *chi.Mux
	listenHttp      net.Listener
	host            string
	port            int32
	shutdownTimeout int16
}

// New return WS object.
func New(opts ...Option) (*WS, error) {
	ws := &WS{
		host:            apps.ReplicaID(),
		router:          chi.NewRouter(),
		port:            8080,
		shutdownTimeout: 10,
	}
	for _, opt := range opts {
		opt(ws)
	}

	listenHttp, err := reuseport.Listen("tcp4", fmt.Sprintf(":%d", ws.port))
	if err != nil {
		return nil, err
	}
	ws.listenHttp = newGracefulListener(listenHttp, time.Second*time.Duration(ws.shutdownTimeout))
	return ws, nil
}

// Run web-server.
func (ws *WS) Run(ctx context.Context, l Logger) {
	if err := ws.validate(); err != nil {
		l.Printf(
			"Run service error",
			map[string]any{
				"error": err.Error(),
			},
		)
		return
	}

	handleShutdown := func(reason string) {
		msg := "no errors"
		if err := ws.shutdown(); err != nil {
			msg = err.Error()
		}
		l.Printf(
			fmt.Sprintf("%s, server shutdown", reason),
			map[string]any{
				"error": msg,
			},
		)
	}

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		mainRouter := chi.NewRouter()
		mainRouter.Mount("/", ws.router)
		mainRouter.Mount("/sys", sysRouter())

		l.Printf(
			"Starting server",
			map[string]any{
				"pid":  os.Getpid(),
				"port": ws.port,
				"host": ws.host,
			},
		)
		errChan <- http.Serve(ws.listenHttp, mainRouter)
	}()
	select {
	case err := <-errChan:
		l.Printf(
			"Server error",
			map[string]any{
				"error": err.Error(),
			},
		)
	case <-ctx.Done():
		handleShutdown("Context canceled")
	case <-sigChan:
		handleShutdown("Received an interrupt signal")
	}
}

// Router object.
func (ws *WS) Router() *chi.Mux {
	return ws.router
}

// Shutdown web-server.
func (ws *WS) shutdown() error {
	return ws.listenHttp.Close()
}

func (ws *WS) validate() error {
	if ws.port > 65535 || ws.port < 1024 {
		return errors.New("incorrect port number, should be in: [1024..65535]")
	}
	return nil
}
