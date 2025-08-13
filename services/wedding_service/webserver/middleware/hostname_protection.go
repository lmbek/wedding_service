package middleware

import (
	"net/http"
)

func ProtectHostNames(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "imbek.dk" && r.Host != "www.imbek.dk" {
			http.Error(w, "Misdirected Request: Host not allowed", http.StatusMisdirectedRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
