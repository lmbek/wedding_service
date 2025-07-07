package website

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"wedding_service/webserver/website/frontend"
)

func loadLayout(w http.ResponseWriter, tmpl *template.Template) {
	file, _ := frontend.DefaultFrontend.GetPrivateFileSystem().Open("layouts/bryllup.gohtml")
	defer file.Close()

	var buffer bytes.Buffer
	_, _ = buffer.ReadFrom(file)
	_, err := tmpl.Parse(buffer.String())
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}
