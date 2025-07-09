package website

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func loadPage(w http.ResponseWriter, tmpl *template.Template, page string) {
	_, err := tmpl.Parse(page)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}

func executePage(w http.ResponseWriter, render *Render, pagePath string, data any) {
	tmpl := template.New(pagePath)

	page := fetchPage(w, render, tmpl, pagePath)
	if page == "" {
		return
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
	}
}

// fetchPage returns a page based on the path given, if it exists as a file
func fetchPage(w http.ResponseWriter, render *Render, tmpl *template.Template, pagePath string) string {
	// Check cache first
	content, exists := ReadCache(pagePath)
	if exists {
		loadLayout(w, render, tmpl)
		loadComponents(w, render, tmpl)

		_, err := tmpl.New("page").Parse(content)
		if err != nil {
			http.Error(w, "Failed to parse cached page template", http.StatusInternalServerError)
			return ""
		}
		return content
	}

	// Not in cache: read from FS
	loadLayout(w, render, tmpl)
	loadComponents(w, render, tmpl)

	filesystem := render.frontend.GetPrivateFileSystem()
	file, err := filesystem.Open(pagePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template open failed: %s — %v", pagePath, err), http.StatusInternalServerError)
		return ""
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(file); err != nil {
		http.Error(w, fmt.Sprintf("Template read failed: %s — %v", pagePath, err), http.StatusInternalServerError)
		return ""
	}
	page := buf.String()

	_, err = tmpl.New("page").Parse(page)
	if err != nil {
		http.Error(w, "Failed to parse page template", http.StatusInternalServerError)
		return ""
	}

	UpdateCache(pagePath, page)
	return page
}
