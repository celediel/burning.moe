package models

import "html/template"

// Link holds an http link, text to display, and an iconify icon class
type Link struct {
	Href, Icon template.URL
	Text       string
}

// TemplateData holds data sent from handlers to templates.
type TemplateData struct {
	StringMap map[string]string      `json:"StringMap" yaml:"StringMap" toml:"StringMap"`
	IntMap    map[string]int         `json:"IntMap" yaml:"IntMap" toml:"IntMap"`
	FloatMap  map[string]float32     `json:"FloatMap" yaml:"FloatMap" toml:"FloatMap"`
	LinkMap   map[string][]Link      `json:"LinkMap" yaml:"LinkMap" toml:"LinkMap"`
	Data      map[string]interface{} `json:"Data" yaml:"Data" toml:"Data"`
	CSRFToken string                 `json:"Csrftoken" yaml:"Csrftoken" toml:"Csrftoken"`
	Flash     string                 `json:"Flash" yaml:"Flash" toml:"Flash"`
	Warning   string                 `json:"Warning" yaml:"Warning" toml:"Warning"`
	Error     string                 `json:"Error" yaml:"Error" toml:"Error"`
}
