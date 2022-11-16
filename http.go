package kocto

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Envelop = map[string]any

func HealthHandler(srv *echo.Echo, cfg Config) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(200, Envelop{
			"status": "available",
			"config": cfg,
			"routes": srv.Routes(),
		})
	}
}

func DefaultMiddleware(srv *echo.Echo, l Logger) {
	srv.Use(middleware.CORS())
	srv.Use(LogMiddleware(l))
}

func LogMiddleware(l Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRoutePath:     true,
		LogURI:           true,
		LogMethod:        true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogQueryParams:   []string{"code", "state", "id", "redirect_uri"},
		LogLatency:       true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			msg := "request"
			values := []any{
				"URI", v.RoutePath,
				"method", v.Method,
				"status", v.Status,
				"error", v.Error,
				"length", v.ContentLength,
				"query", queryString(v.QueryParams),
				"latency", v.Latency.Milliseconds(),
			}

			if v.RoutePath == "/healthz" {
				return nil
			}

			if v.Status >= 400 {
				l.Errorw(msg, values...)
			} else {
				l.Infow(msg, values...)
			}

			return nil
		},
	})
}

func queryString(q map[string][]string) string {
	qString := "?"

	for query, values := range q {
		vString := ""

		for _, val := range values {
			if vString != "" {
				vString += ","
			}

			vString += val
		}

		if qString != "?" {
			qString += "&"
		}

		qString += query + "=" + vString
	}

	return qString
}

func RunServer(srv *echo.Echo, cfg Config, l Logger) error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		l.Infow("shutting down http server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// If there is need for background tasks started from http requests
		// wait here (sync.WaitGroup)
		// a.wg.Wait()
		shutdownError <- nil
	}()

	l.Info("starting http server on port ", cfg.Port)

	err := srv.Start(":" + cfg.Port)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	return nil
}
