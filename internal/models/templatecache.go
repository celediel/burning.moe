package models

import (
	"bytes"
	"html/template"
	"net/http"
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

// Execute writes the template to the supplied writer using the supplied data.
func (self *TemplateCacheItem) Execute(d *TemplateData, w http.ResponseWriter) error {
	buf := new(bytes.Buffer)
	err := self.Template.Execute(buf, d)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
