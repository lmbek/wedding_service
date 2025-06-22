package website

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed frontend/out/private/components/*.gohtml
var componentsData embed.FS

func LoadComponents(w http.ResponseWriter, tmpl *template.Template) {
	_, err := tmpl.ParseFS(componentsData, "frontend/out/private/components/*.gohtml")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}
