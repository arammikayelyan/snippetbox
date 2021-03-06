package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/arammikayelyan/snippetbox/pkg/forms"
	"github.com/arammikayelyan/snippetbox/pkg/models"
)

// templateData type is a holding structure for
// any dynamic data that is passing to html template
type templateData struct {
	AuthenticatedUser *models.User
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
}

// humanDate returns a nicely formatted string representation of a time.Time object
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global vatiable.
// This is essentially a string-keyed map which acts as a lookup between
// the names of object's custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		// The template.FuncMap must be registered with the template set
		// before call the ParseFiles() method. This means we have to use
		// template.New() to create an empty template set, use the Funcs()
		// method to register the template.FuncMap, and then parse the file
		// as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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
