package main

import (
	"net/http"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routes handles all of the HTTP setup. Middleware is enabled,
// static fileserver is setup, and handlers are ... handled
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Import some middleware
	mux.Use(middleware.Recoverer)

	// Setup static file server
	app.Logger.Debug("Setting up /static file server")
	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// Setup routes for handlers
	for _, handler := range handlers.Handlers {
		app.Logger.Info("Setting up handler for " + handler.Handles)
		mux.Get(handler.Handles, handler.Handler)
	}

	// Setup home handler
	mux.Get("/", handlers.HomeHandler)

	return mux
}
