package handlers

import (
	"net/http"

	"git.burning.moe/celediel/burning.moe/internal/config"
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

// makeBasicHandler returns a simple handler that renders a template from `name`.page.tmpl
func makeBasicHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Infof("Got request for %s page", name)
		pageName := name + ".page.tmpl"
		render.RenderTemplate(w, pageName)
	}
}
