package website

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed frontend/out/private/layouts/bryllup.gohtml
var layoutData string

func LoadLayout(w http.ResponseWriter, tmpl *template.Template) {
	_, err := tmpl.Parse(layoutData)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}
