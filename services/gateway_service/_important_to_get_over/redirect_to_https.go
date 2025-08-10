package webserver

import (
	"fmt"
	"net"
	"net/http"
)

func redirectToHTTPS(httpsPort string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		// Split host and port (if any)
		hostName, _, err := net.SplitHostPort(host)
		if err != nil {
			// no port specified, just append the httpsPort if not default 443
			hostName = host
		}

		// Replace port with httpsPort unless httpsPort is "443"
		if httpsPort == "443" {
			host = hostName // no port in host for standard HTTPS
		} else {
			host = net.JoinHostPort(hostName, httpsPort)
		}

		target := "https://" + host + r.URL.RequestURI()
		fmt.Println(target)
		http.Redirect(w, r, target, http.StatusPermanentRedirect)
	})
}
