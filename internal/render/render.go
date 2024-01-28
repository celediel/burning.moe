package render

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/models"
	"git.burning.moe/celediel/burning.moe/internal/td"
)

const (
	templatesDir string = "./templates/"
	layoutGlob   string = "*.layout.tmpl"
	pageGlob     string = "*.tmpl"
)

var app *config.AppConfig

// Initialise the render package.
func Initialise(a *config.AppConfig) {
	app = a
	if app.UseCache {
		var err error
		app.TemplateCache, err = generateNewTemplateCache()
		if err != nil {
			app.Logger.Fatal("Error generating template cache, bailing out!")
		}
	}
}

// RenderTemplate renders requested template (t), pulling from cache.
func RenderTemplate(w http.ResponseWriter, filename string) {
	// TODO: implement this better
	if !app.UseCache {
		regenerateTemplateCache()
	}

	template, err := GetTemplateFromCache(filename)
	if err != nil {
		app.Logger.Fatalf("Tried loading %s from the cache, but %s!", filename, err)
	}

	data, err := GetOrGenerateTemplateData(filename)
	if err != nil {
		app.Logger.Error(err)
	}

	app.Logger.Debug(fmt.Sprintf("Executing template %s", filename), "data", &data)
	err = template.Execute(data, w)
	if err != nil {
		app.Logger.Fatalf("Failed to execute template %s: %s", filename, err)
	}
}

func RenderTemplateWithData(w http.ResponseWriter, filename string, data *models.TemplateData) {
	if !app.UseCache {
		regenerateTemplateCache()
	}

	template, err := GetTemplateFromCache(filename)
	if err != nil {
		app.Logger.Fatalf("Tried loading %s from the cache, but %s!", filename, err)
	}

	app.Logger.Debug(fmt.Sprintf("Executing template %s", filename), "data", &data)
	err = template.Execute(data, w)
	if err != nil {
		app.Logger.Fatalf("Failed to execute template %s: %s", filename, err)
	}
}

// GetTemplateFromCache gets templates from cache
func GetTemplateFromCache(filename string) (*models.TemplateCacheItem, error) {
	if template, ok := app.TemplateCache.Cache[filename]; ok {
		return &template, nil
	} else {
		return &models.TemplateCacheItem{}, errors.New("Couldn't load template from cache")
	}
}

// GetOrGenerateTemplateData gets template data from file, or generate simple
func GetOrGenerateTemplateData(filename string) (*models.TemplateData, error) {
	template, err := GetTemplateFromCache(filename)
	if err != nil {
		return &models.TemplateData{}, err
	}

	data, err := td.LoadTemplateData(filename)
	if err == nil {
		app.Logger.Debug(fmt.Sprintf("Loaded data for template %s.", filename), "data", &data)
		if _, ok := data.StringMap["GeneratedAt"]; !ok {
			data.StringMap["GeneratedAt"] = template.GeneratedAt.Format(time.UnixDate)
		}
	} else {
		app.Logger.Info(fmt.Sprintf("Loading template data for %s failed, using default template data.", filename), "err", err)
		data = td.MakeBasicTemplateData(template.GeneratedAt)
	}

	return &data, nil
}

// regenerateTemplateCache regenerates the template cache
func regenerateTemplateCache() {
	c, err := generateNewTemplateCache()
	if err != nil {
		app.Logger.Fatal("Error generating template cache, bailing out!")
	}
	app.TemplateCache = c

}

// generateNewTemplateCache generates a new template cache.
func generateNewTemplateCache() (models.TemplateCache, error) {
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
