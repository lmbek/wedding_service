package webserver

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"wedding_service/config"
)

var allowedHostsMap map[string]struct{}

func initAllowedHosts(config config.Config) {
	raw := config.Hostnames() // map[string][]string
	allowedHostsMap = make(map[string]struct{})

	// Populate the allowed hosts map from env config
	for domain, aliases := range raw {
		allowedHostsMap[domain] = struct{}{}
		for _, alias := range aliases {
			allowedHostsMap[alias] = struct{}{}
		}
	}
}

// ProtectHostsMiddleware checks the Host header and rejects requests to disallowed hosts.
func ProtectHostsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for WebSocket upgrade requests and allow them to pass through
		if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
			// Skip host validation for WebSocket connections
			next.ServeHTTP(w, r)
			return
		}

		// Get the host from the request (can have port)
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}

		// Check if the host is allowed
		if _, allowed := allowedHostsMap[host]; !allowed {
			// Block the connection if the host is not allowed
			fmt.Printf("🔒 Blocking connection from unauthorized host: %s\n", host)

			// Fallback to returning Forbidden error if hijacking fails
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Proceed if host is allowed
		next.ServeHTTP(w, r)
	})
}
