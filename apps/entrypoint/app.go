package entrypoint

import (
	"context"
	"strconv"

	"github.com/Mad-Pixels/golang-playground/apps/pkg/log"
	"github.com/Mad-Pixels/golang-playground/apps/pkg/ws"

	"github.com/go-chi/chi/v5"
)

type webServer interface {
	Run(ctx context.Context, l ws.Logger)
	Router() *chi.Mux
}

// App object.
type App struct {
	logger ws.Logger
	server webServer
}

// New create Service object.
func New(listenPort, logLevel string) (App, error) {
	l, err := log.New([]byte(logLevel))
	if err != nil {
		return App{}, err
	}
	p, err := strconv.Atoi(listenPort)
	if err != nil {
		return App{}, err
	}

	s, err := ws.New(
		ws.WithMiddlewareReqID(),
		ws.WithMiddlewareRealIP(),

		ws.WithCustomPort(p),
		ws.WithShutdownTimeout(10),
		ws.WithMiddlewareCustomAfter(customLogger(l)),
	)
	return App{
		server: s,
		logger: l,
	}, nil
}

// Run application.
func (a App) Run(ctx context.Context) {
	a.routerViews()
	a.routerTools()
	a.server.Run(ctx, a.logger)
}
