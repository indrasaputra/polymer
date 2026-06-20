package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/validator"
	wmid "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/http/middleware"
)

// Server holds server data.
type Server struct {
	*echo.Echo
	Port string
}

// New creates an instance of Server with all necessary middleware ready.
func New(cfg *config.Config) (*Server, error) {
	e := echo.New()

	e.Validator = validator.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.ContextTimeout(time.Duration(cfg.GlobalTimeoutInSeconds) * time.Second))
	e.Use(middleware.RequestID())

	// TODO: move to sdk logger and middleware
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogLatency:   true,
		LogProtocol:  true,
		LogRemoteIP:  true,
		LogRequestID: true,
		LogURI:       true,
		LogRoutePath: true,
		HandleError:  true, // forwards the error to the global error handler so it can pick the status code
		LogValuesFunc: func(_ *echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	jwtmid, err := wmid.NewJwtMiddleware(cfg)
	if err != nil {
		return nil, err
	}
	e.Use(jwtmid)

	// TODO: open telemetry middleware

	return &Server{Echo: e, Port: cfg.Port}, nil
}

// StartWithGracefulStop starts the server with graceful stop ready.
func (s *Server) StartWithGracefulStop(ctx context.Context, cfg *config.Config) error {
	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%s", s.Port),
		GracefulTimeout: time.Duration(cfg.GracefulTimeoutInSeconds) * time.Second,
	}

	return sc.Start(ctx, s.Echo)
}

// PrepareForGracefulStop prepares the context and cancel func to be used in StartWithGracefulStop.
// the cancel func must be called in defer.
func (s *Server) PrepareForGracefulStop() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
