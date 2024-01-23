package handlers

import (
	"net/http"

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
		Handler: func(w http.ResponseWriter, r *http.Request) {
			app.Logger.Info("Got request for about page.")
			render.RenderTemplate(w, "about.page", &models.TemplateData{})
		},
	},
	// / comes last
	{
		Handles: "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			app.Logger.Info("Got request for homepage.")
			render.RenderTemplate(w, "home.page", &models.TemplateData{})
		},
	},
}

// Initialise the handlers package.
func Initialise(a *config.AppConfig) {
	app = a
}
