package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/handlers"
)

// routes handles all of the HTTP setup. Middleware is enabled,
// static fileserver is setup, and handlers are ... handled
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Import some middleware
	for _, mw := range Middleware {
		mux.Use(mw)
	}

	// Setup static file server
	app.Logger.Debug("Setting up /static file server")
	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// Setup routes for handlers
	for _, handler := range handlers.Handlers {
		app.Logger.Info("Setting up handler for " + handler.Handles)
		mux.Get(handler.Handles, handler.Handler)
	}

	// Setup extra handlers
	mux.Get("/", handlers.HomeHandler)
	mux.Get("/robots.txt", handlers.RobotHandler)

	return mux
}
