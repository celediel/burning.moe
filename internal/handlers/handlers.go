package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"git.burning.moe/celediel/burning.moe/internal/config"
	"git.burning.moe/celediel/burning.moe/internal/models"
	"git.burning.moe/celediel/burning.moe/internal/render"
	"git.burning.moe/celediel/burning.moe/internal/td"
)

// Handler holds data required for handlers.
type Handler struct {
	Handles string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var app *config.AppConfig

// The actual handlers
var Handlers = []Handler{
	{
		Handles: "/about",
		Handler: makeBasicHandler("about"),
	},
	{
		Handles: "/projects",
		Handler: makeLinksHandler("projects"),
	},
	{
		Handles: "/apps",
		Handler: makeLinksHandler("apps"),
	},
}

// Initialise the handlers package.
func Initialise(a *config.AppConfig) {
	app = a
}

// HomeHandler handles /, generating data from Handlers
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Info("Got request for homepage")

	page := "home.page.tmpl"
	d := models.TemplateData{}

	t, err := render.GetTemplateFromCache(page)
	if err != nil {
		app.Logger.Error(fmt.Sprintf("couldn't get %s from cache", page), "err", err)
		d.StringMap = map[string]string{
			"GeneratedAt": time.Now().Format(time.UnixDate),
		}
	} else {
		d.StringMap = map[string]string{
			"GeneratedAt": t.GeneratedAt.Format(time.UnixDate),
		}
	}

	var pages []models.Link = []models.Link{}

	for _, handler := range Handlers {
		href := strings.TrimPrefix(handler.Handles, "/")
		pages = append(pages, models.Link{
			Href: template.URL(href),
			Text: href,
		})
	}

	d.LinkMap = make(map[string][]models.Link)
	d.LinkMap["Pages"] = pages
	app.Logger.Debug("handling home with some data", "data", &d)
	render.RenderTemplateWithData(w, "home.page.tmpl", &d)
}

// makeBasicHandler returns a simple handler that renders a template from `name`.page.tmpl
func makeBasicHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Infof("Got request for %s page", name)
		pageName := name + ".page.tmpl"
		render.RenderTemplate(w, pageName)
	}
}

// makeLinksHandler returns a handler for links.tmpl with template data from `name`
func makeLinksHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		page := "links.tmpl"
		template, err := render.GetTemplateFromCache(page)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("couldn't get %s from cache", page), "err", err)
		}

		app.Logger.Infof("Got request for %s links page", name)
		data, err := td.LoadTemplateData(name)
		if err != nil {
			app.Logger.Fatal("couldn't load template data for "+name, "err", err)
		} else {
			data.StringMap["GeneratedAt"] = template.GeneratedAt.Format(time.UnixDate)
		}

		app.Logger.Debug("handling a links page", "data", &data)
		render.RenderTemplateWithData(w, page, &data)
	}
}
