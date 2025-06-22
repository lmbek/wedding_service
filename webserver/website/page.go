package website

import (
	"fmt"
	"html/template"
	"net/http"
)

func LoadPage(w http.ResponseWriter, frontpageData string, tmpl *template.Template) {
	LoadLayout(w, tmpl)
	LoadComponents(w, tmpl)
	_, err := tmpl.Parse(frontpageData)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}

func ExecutePage(w http.ResponseWriter, tmpl *template.Template, frontPageData string, data any) {
	LoadPage(w, frontPageData, tmpl)

	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}
