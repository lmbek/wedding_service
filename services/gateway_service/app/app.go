package app

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gateway_service/certificate"
	"gateway_service/config"
	"gateway_service/logger"
	"gateway_service/middleware"
	"gateway_service/reverseproxy"
	"gateway_service/tracer"
)

type App interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type app struct {
	ctx          context.Context
	cancel       context.CancelCauseFunc
	config       config.Config
	logger       logger.Logger
	tracer       tracer.Tracer
	acme         certificate.AutoCertManager
	reverseProxy reverseproxy.ReverseProxy
}

func NewApp(cfg config.Config, log logger.Logger, tracer tracer.Tracer, acme certificate.AutoCertManager) (App, error) {
	ctx, cancel := context.WithCancelCause(context.Background())

	app := &app{
		ctx:    ctx,
		cancel: cancel,
		config: cfg,
		logger: log,
		tracer: tracer,
		acme:   acme,
	}

	return app, app.init()
}

func (a *app) init() error {
	// Init outbound transport and middleware
	outboundTransport := a.buildOutboundTransport()
	newMiddleware := middleware.NewMiddleware(a.logger, a.tracer)

	// Build HTTP :80 server (ACME HTTP-01 + redirect)
	httpServer := reverseproxy.NewHTTPServer(reverseproxy.HTTPServer{
		Port: a.config.HTTPPort(),
		ACME: a.acme,
		RedirectFunc: func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		},
	})

	// Build HTTPS-server :443
	httpsServer, err := reverseproxy.NewHTTPSServer(reverseproxy.HTTPSServer{
		Port:           a.config.HTTPSPort(),
		SelfSignedCert: a.config.CertPath(), // brug dine egne stier
		SelfSignedKey:  a.config.KeyPath(),
		ACME:           a.acme,
		// Handler is set later with a.reverseProxy.Handler()
		PreferHTTP2Only: true,
	})
	if err != nil {
		return fmt.Errorf("failed to build HTTPS server: %w", err)
	}

	// Opret reverse proxy med servere
	reverseProxy := reverseproxy.NewReverseProxy(
		a.config,
		httpServer,
		httpsServer,
		newMiddleware,
		outboundTransport,
	)

	// Gem reverse proxy på app før brug
	a.reverseProxy = reverseProxy

	// Konfigurer backends
	backends := a.config.Backends()
	for host, target := range backends {
		link, err := url.Parse(target)
		if err != nil {
			return err
		}
		err = a.reverseProxy.AddHost(host, link)
		if err != nil {
			return err
		}
		// Optional: also register www.<host> if missing (production-friendly default)
		if !strings.HasPrefix(host, "www.") {
			www := "www." + host
			if _, ok := backends[www]; !ok {
				_ = a.reverseProxy.AddHost(www, link)
			}
		}
		// Warmup: best-effort reachability check (TCP dial with timeout) for diagnostics
		hostport := link.Host
		if link.Scheme == "http" {
			if !strings.Contains(hostport, ":") {
				hostport = hostport + ":80"
			}
		}
		if link.Scheme == "https" {
			if !strings.Contains(hostport, ":") {
				hostport = hostport + ":443"
			}
		}
		d := net.Dialer{Timeout: 2 * time.Second}
		conn, dialErr := d.DialContext(a.ctx, "tcp", hostport)
		if dialErr == nil {
			_ = conn.Close()
		} else {
			// Do not fail startup; just log via logger if available (keeps minimal impact)
			// This helps explain potential 502 if backend is unreachable inside the network
		}
	}

	// Sæt HTTPS handler nu hvor reverse proxy er klar
	httpsServer.Handler = a.reverseProxy.Handler()

	// Tilpas HTTP :80 handleren til at tillade WebSocket-upgrade anmodninger
	// gennem reverse proxy, mens alm. HTTP stadig redirectes til HTTPS og ACME
	// udfordringer håndteres.
	httpServer.Handler = reverseproxy.NewHTTPRedirectOrProxyHandler(
		a.acme,
		func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		},
		a.reverseProxy.Handler(),
	)

	return nil
}

func (a *app) Run() error {
	err := a.reverseProxy.ListenAndServe()
	if err != nil {
		return err
	}
	defer a.reverseProxy.Close()

	return nil
}

func (a *app) Shutdown(ctx context.Context) error {
	var err, httpErr, httpsErr error

	if a.reverseProxy.HTTPServer() != nil {
		httpsErr = a.reverseProxy.HTTPServer().Shutdown(ctx)
	}

	if a.reverseProxy.HTTPSServer() != nil {
		httpErr = a.reverseProxy.HTTPSServer().Shutdown(ctx)
	}

	err = ctx.Err()
	a.cancel(err)
	return errors.Join(err, httpErr, httpsErr)
}

func (a *app) buildOutboundTransport() http.RoundTripper {
	base := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        200,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return a.tracer.NewTransport(base)
}
