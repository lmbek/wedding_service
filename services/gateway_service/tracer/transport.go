package tracer

import "net/http"

type Transport interface {
	http.RoundTripper
}

type transport struct {
	roundTripper http.RoundTripper
	tracer       *tracer
}

func (t *tracer) NewTransport(roundTripper http.RoundTripper) Transport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	return &transport{roundTripper: roundTripper, tracer: t}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	id, ok := t.tracer.GetID(req.Context())
	if ok && id != "" {
		clone := req.Clone(req.Context())
		clone.Header.Set(t.tracer.headerKey, id)
		return t.roundTripper.RoundTrip(clone)
	}
	return t.roundTripper.RoundTrip(req)
}
