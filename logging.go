package kocto

import (
	"github.com/labstack/echo"
	"github.com/mbcrocci/pika"
)

// WithRequest extracts information from the request context and adds it to the log context
func WithRequest(l Logger, c echo.Context) Logger {
    return l.With(
		"URI", c.Path(),
		"query", c.QueryString(),
	)
}

// WithEventContext extracts information from the consumerOptions and adds to the log context
func WithEventContext(l Logger, o pika.ConsumerOptions) Logger {
	return l.With(
		"topic", o.Topic,
	)
}
