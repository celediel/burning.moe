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
func RenderTemplate(w http.ResponseWriter, filename string) {
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
		app.Logger.Errorf("Couldn't get %s from template cache, dunno what happened, but we're gonna generate a new one", filename)
		c, err := GenerateNewTemplateCache()
		if err != nil {
			app.Logger.Fatal("Error generating template cache, bailing out!")
		}
		app.TemplateCache = c
		template = app.TemplateCache.Cache[filename]
	}

	// Get template data from file, or generate simple
	data, err := models.LoadTemplateData(filename)
	if err == nil {
		app.Logger.Debug(fmt.Sprintf("Loaded data for template %s.", filename), "data", data)
		if _, ok := data.StringMap["GeneratedAt"]; !ok {
			data.StringMap["GeneratedAt"] = template.GeneratedAt.Format(time.UnixDate)
		}
	} else {
		app.Logger.Info(fmt.Sprintf("Loading template data for %s failed, using default template data.", filename), "err", err)
		data = models.MakeBasicTemplateData(template.GeneratedAt)
	}

	// Execute templates in a new buffer
	buf := new(bytes.Buffer)
	err = template.Template.Execute(buf, data)
	if err != nil {
		app.Logger.Fatal(fmt.Sprintf("Error executing template %s! Goodbye!", filename), "err", err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		app.Logger.Error(fmt.Sprintf("Error writing template %s!", filename), "err", err)
	}
}
