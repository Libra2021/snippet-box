package main

import (
	"snippetbox.libra.dev/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
