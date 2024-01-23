package models

import (
	"html/template"
	"time"
)

// TemplateCache holds the template cache as map of TemplateCacheItem
type TemplateCache struct {
	Cache map[string]TemplateCacheItem
}

// TemplateCacheItem holds a pointer to a generated
// template, and the time it was generated at.
type TemplateCacheItem struct {
	Template    *template.Template
	GeneratedAt time.Time
}
