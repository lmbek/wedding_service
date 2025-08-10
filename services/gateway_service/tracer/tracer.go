package tracer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Tracer defines the tracing API.
// Interface-first with unexported concrete implementation.

type Tracer interface {
	NewID() string
	WithID(ctx context.Context, id string) context.Context
	GetID(ctx context.Context) (string, bool)
	Middleware(next http.Handler) http.Handler
	NewTransport(roundTripper http.RoundTripper) Transport
}

type tracer struct {
	key       contextKey
	headerKey string
}

// NewTracer returns a Tracer without exposing the concrete struct.
func NewTracer(header string) Tracer {
	if header == "" {
		header = "X-Trace-Id"
	}
	return &tracer{key: contextKey{}, headerKey: header}
}

func (t *tracer) NewID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func (t *tracer) WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, t.key, id)
}

func (t *tracer) GetID(ctx context.Context) (string, bool) {
	v := ctx.Value(t.key)
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

func (t *tracer) Middleware(next http.Handler) http.Handler {
	if next == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(t.headerKey)
			if id == "" {
				id = t.NewID()
			}
			w.Header().Set(t.headerKey, id)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(t.headerKey)
		if id == "" {
			id = t.NewID()
		}
		ctx := t.WithID(r.Context(), id)
		w.Header().Set(t.headerKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
