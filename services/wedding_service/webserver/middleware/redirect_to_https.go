package middleware

import "net/http"

// RedirectToHTTPS Redirect HTTP traffic on port 80 to HTTPS
func RedirectToHTTPS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := "https://" + r.Host + r.URL.RequestURI()
		http.Redirect(w, r, target, http.StatusPermanentRedirect)
	})
}
