package website

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

// loadComponents parses all *.gohtml in the components dir into baseTmpl.
func loadComponents(w http.ResponseWriter, render *Render, baseTmpl *template.Template) {

	privateFS := render.frontend.PrivateFS()

	err := fs.WalkDir(privateFS, "components", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".gohtml") {
			return nil
		}

		content, err := fs.ReadFile(privateFS, path)
		if err != nil {
			return fmt.Errorf("failed to read component %q: %w", path, err)
		}

		// Register component with its full relative path
		_, err = baseTmpl.New(path).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse component %q: %w", path, err)
		}
		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading components: %v", err), http.StatusInternalServerError)
	}
}
