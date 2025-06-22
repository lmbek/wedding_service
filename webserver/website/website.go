package website

import "net/http"

func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write([]byte("kage<style>body { color: red; }</style>"))
}

func InvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write([]byte("Invitation<style>body { color: pink; }</style>"))
}
