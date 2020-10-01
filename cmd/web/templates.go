package main

import (
	"path/filepath"
	"text/template"

	"github.com/arammikayelyan/snippetbox/pkg/models"
)

// templateData type is a holding structure for
// any dynamic data that is passing to html template
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// newTemplateCache caches the templates in memory instead of disk
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the page one-by-one
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full path
		name := filepath.Base(page)

		// Parse the page files into a template set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob to add any layout templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob to add any partial templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		cache[name] = ts
	}

	return cache, nil
}
