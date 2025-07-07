package website

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"wedding_service/webserver/website/frontend"
)

func loadLayout(w http.ResponseWriter, tmpl *template.Template) {
	file, err := frontend.DefaultFrontend.GetPrivateFileSystem().Open("layouts/bryllup.gohtml")
	if err != nil {
		http.Error(w, "failed to open layout: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(file)
	if err != nil {
		http.Error(w, "failed to read layout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tmpl.Parse(buffer.String())
	if err != nil {
		http.Error(w, "failed to parse layout: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
