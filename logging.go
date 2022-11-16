package kocto

import (
	"github.com/labstack/echo"
	"github.com/mbcrocci/pika"
)

func (l *Logger) WithRequest(c echo.Context) *Logger {
	l.SugaredLogger = l.SugaredLogger.With(
		"URI", c.Path(),
		"query", c.QueryString(),
	)

	return l
}

func (l *Logger) WithEventContext(o pika.ConsumerOptions) *Logger {
	l.SugaredLogger = l.SugaredLogger.With(
		"topic", o.Topic,
	)

	return l
}
