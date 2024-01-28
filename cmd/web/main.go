// Main entry point for the web app. Does all
// the setup, then runs the http server.
package main

import (
	"fmt"
	"net/http"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/handlers"
	"git.burning.moe/celediel/burning.moe/internal/render"
)

// App wide config data and such
var app config.AppConfig

func main() {
	// Initialise app and config
	app = config.Initialise()

	// Initialise handlers and renderer
	handlers.Initialise(&app)
	render.Initialise(&app)

	// Initialise the webserver
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.ListenPort),
		Handler: routes(&app),
	}

	// and finally, start the server
	app.Logger.Printf("Starting HTTP Server on port %d", app.ListenPort)
	err := srv.ListenAndServe()
	if err != nil {
		app.Logger.Fatal("Failed to start HTTP Server!", "err", err)
	}
}
