package main

import (
	"net/http"

	"github.com/arl/statsviz"
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

	// Setup stats viewer
	stats, _ := statsviz.NewServer()
	mux.Get("/debug/statsviz/ws", stats.Ws())
	mux.Get("/debug/statsviz", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/debug/statsviz/", 301)
	})
	mux.Handle("/debug/statsviz/*", stats.Index())

	// Setup static file server
	app.Logger.Debug("Setting up /static file server")
	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// handle certbot challenge
	mux.Handle("/.well-known/acme-challenge/*", http.StripPrefix("/.well-known/acme-challenge", http.FileServer(http.Dir("./.well-known/acme-challenge"))))

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
