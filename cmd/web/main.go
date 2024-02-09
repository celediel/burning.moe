// Main entry point for the web app. Does all
// the setup, then runs the http server.
package main

import (
	"fmt"
	"net/http"
	"time"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/handlers"
	"git.burning.moe/celediel/burning.moe/internal/render"

	"github.com/madflojo/tasks"
)

// App wide config data and such
var app config.AppConfig

func main() {
	// Initialise app and config
	app = config.Initialise()

	// Initialise handlers and renderer
	handlers.Initialise(&app)
	render.Initialise(&app)

	// Setup task to regenerate template cache
	scheduler := tasks.New()
	defer scheduler.Stop()
	id, err := scheduler.Add(&tasks.Task{
		Interval: time.Duration(app.CacheTimer),
		TaskFunc: func() error {
			err := render.RegenerateTemplateCache()
			if err != nil {
				return err
			}
			app.Logger.Infof("Regenerate cache job finished at %s", time.Now())
			return nil
		},
		ErrFunc: func(err error) {
			app.Logger.Error("Error in template regeneration cache job!", "err", err)
		},
	})
	if err != nil {
		app.Logger.Error("Error setting upscheduler", "err", err)
	}

	app.Logger.Info("Started cache regeneration task.", "interval", app.CacheTimer.String(), "id", id)

	// Initialise the webserver
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.ListenPort),
		Handler: routes(&app),
	}

	// and finally, start the server
	app.Logger.Printf("Starting HTTP Server on port %d", app.ListenPort)
	err = srv.ListenAndServe()
	if err != nil {
		app.Logger.Fatal("Failed to start HTTP Server!", "err", err)
	}
}
