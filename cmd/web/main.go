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

func main() {
	// Initialise app and config
	app := config.Initialise()

	// Initialise handlers and renderer
	handlers.Initialise(&app)
	render.Initialise(&app)

	// Initialise the webserver
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.ListenPort),
		Handler: routes(&app),
	}

	// and finally, start the server
	app.Logger.Printf("Listening on port %d", app.ListenPort)
	srv.ListenAndServe()
}