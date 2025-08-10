package middleware

import (
	"gateway_service/logger"
	"gateway_service/tracer"
	"net/http"
)

type Middleware interface {
	Handler(next http.Handler) http.Handler
}

type middleware struct {
	tracer tracer.Tracer
	logger logger.Logger
}

func NewMiddleware(logger logger.Logger, t tracer.Tracer) Middleware {
	return &middleware{tracer: t}
}

func (m *middleware) Handler(next http.Handler) http.Handler {
	return m.tracer.Middleware(next)
}
