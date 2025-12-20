package main

import (
	"errors"

	"snippetbox.kira.net/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

// !IMPORTANT: This a little deviation from the author's implementation
func NewTemplateDatum(s *models.Snippet) (*templateData, error) {
	if s == nil {
		return nil, errors.New("templates: Pointer *models.Snippet is nil")
	}
	return &templateData{
		Snippet: *s,
	}, nil
}

func NewTemplateData(s *[]models.Snippet) (*templateData, error) {
	if len(*s) == 0 {
		return nil, errors.New("templates: Pointer *models.Snippet is nil")
	}
	return &templateData{
		Snippets: *s,
	}, nil
}
