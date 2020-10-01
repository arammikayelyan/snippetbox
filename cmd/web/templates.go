package main

import "github.com/arammikayelyan/snippetbox/pkg/models"

// templateData type is a holding structure for
// any dynamic data that is passing to html template
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
