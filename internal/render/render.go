package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/models"
)

const (
	templatesDir string = "./templates/"
	layoutGlob   string = "*.layout.tmpl"
	pageGlob     string = "*.page.tmpl"
)

var app *config.AppConfig

// Initialise the render package.
func Initialise(a *config.AppConfig) {
	var err error
	app = a
	if app.UseCache {
		app.TemplateCache, err = GenerateNewTemplateCache()
	}
	if err != nil {
		app.Logger.Fatal("Error generating template cache, bailing out!")
	}
}

// GenerateNewTemplateCache generates a new template cache.
func GenerateNewTemplateCache() (models.TemplateCache, error) {
	// start with an empty map
	cache := models.TemplateCache{}
	cache.Cache = map[string]models.TemplateCacheItem{}

	// Generate a list of pages based on globs
	pages, err := filepath.Glob(templatesDir + pageGlob)

	// a nice try catch would be pretty cool right about here
	if err != nil {
		return cache, err
	}

	// Iterate each page, parsing the file and adding it to the cache
	for _, page := range pages {
		name := filepath.Base(page)
		app.Logger.Info("Generating template " + name)
		generatedAt := time.Now()

		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		// Glob and parse any layouts found
		layouts, err := filepath.Glob(templatesDir + layoutGlob)
		if err != nil {
			return cache, err
		}

		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob(templatesDir + layoutGlob)
			if err != nil {
				return cache, err
			}
		}
		cache.Cache[name] = models.TemplateCacheItem{
			Template:    templateSet,
			GeneratedAt: generatedAt,
		}
		app.Logger.Debugf("Generated %s at %v", name, generatedAt.Format(time.UnixDate))
	}

	// All was good, so return the cache, and no error
	return cache, nil
}

// RenderTemplate renders requested template (t), pulling from cache.
func RenderTemplate(w http.ResponseWriter, filename string, data *models.TemplateData) {
	if !app.UseCache {
		c, err := GenerateNewTemplateCache()
		if err != nil {
			app.Logger.Fatal("Error generating template cache, bailing out!")
		}
		app.TemplateCache = c
	}

	// Get templates from cache
	template, ok := app.TemplateCache.Cache[filename]
	if !ok {
		app.Logger.Fatal(fmt.Sprintf("Couldn't get %s from template cache, bailing out!", filename))
	}

	// Execute templates in a new buffer
	buf := new(bytes.Buffer)
	err := template.Template.Execute(buf, data)

	if err != nil {
		app.Logger.Fatal(fmt.Sprintf("Error executing template %s! Goodbye!", filename), "err", err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		app.Logger.Error(fmt.Sprintf("Error writing template %s!", filename), "err", err)
	}
}
