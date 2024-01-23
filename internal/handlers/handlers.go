package handlers

import (
	"net/http"
	"time"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/models"
	"git.burning.moe/celediel/burning.moe/internal/render"
)

// Handler holds data required for handlers.
type Handler struct {
	Handles string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var app *config.AppConfig

// The actual handlers
var Handlers = []Handler{
	// /about
	{
		Handles: "/about",
		Handler: makeBasicHandler("about"),
	},
	// / comes last
	{
		Handles: "/",
		Handler: makeBasicHandler("home"),
	},
}

// Initialise the handlers package.
func Initialise(a *config.AppConfig) {
	app = a
}

// makeBasicTemplateData creates a blank TemplateData containing only the
// time the related template was generated
func makeBasicTemplateData(name string) models.TemplateData {
	var strMap map[string]string
	if _, ok := app.TemplateCache.Cache[name]; ok {
		strMap = map[string]string{
			"GeneratedAt": app.TemplateCache.Cache[name].GeneratedAt.Format(time.UnixDate),
		}
	} else {
		strMap = make(map[string]string)
	}

	templateData := models.TemplateData{
		StringMap: strMap,
	}
	return templateData
}

// makeBasicHandler creates a basic handler that builds from a .page.tmpl
// file, and sends only the time the template was generated as TemplateData
func makeBasicHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Infof("Got request for %s page", name)
		pageName := name + ".page.tmpl"
		templateData := makeBasicTemplateData(pageName)
		render.RenderTemplate(w, pageName, &templateData)
	}
}
