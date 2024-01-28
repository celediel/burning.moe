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
	Handler func(writer http.ResponseWriter, request *http.Request)
}

var app *config.AppConfig

// The actual handlers
var Handlers = []Handler{
	{
		Handles: "/about",
		Handler: makeBasicHandler("about"),
	},
	{
		Handles: "/apps",
		Handler: makeLinksHandler("apps"),
	},
	{
		Handles: "/projects",
		Handler: makeLinksHandler("projects"),
	},
}

// Initialise the handlers package.
func Initialise(a *config.AppConfig) {
	app = a
}

// HomeHandler handles /, generating data from Handlers
func HomeHandler(writer http.ResponseWriter, request *http.Request) {
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
	render.RenderTemplateWithData(writer, "home.page.tmpl", &d)
}

// RobotHandler creates a handler for robots.txt out of existing Handlers
func RobotHandler(writer http.ResponseWriter, request *http.Request) {
	var robots strings.Builder
	robots.WriteString(fmt.Sprintf("User-agent: %s\nAllow: /\n", request.UserAgent()))

	for _, handler := range Handlers {
		robots.WriteString(fmt.Sprintf("Allow: %s\n", handler.Handles))
	}

	robots.WriteString("Disallow: /*\n")

	fmt.Fprint(writer, robots.String())
}

// makeBasicHandler returns a simple handler that renders a template from `name`.page.tmpl
func makeBasicHandler(name string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		pageName := name + ".page.tmpl"
		render.RenderTemplate(writer, pageName)
	}
}

// makeLinksHandler returns a handler for links.tmpl with template data from `name`
func makeLinksHandler(name string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		page := "links.tmpl"
		template, err := render.GetTemplateFromCache(page)
		if err != nil {
			app.Logger.Error(fmt.Sprintf("couldn't get %s from cache", page), "err", err)
		}

		data, err := td.LoadTemplateData(name)
		if err != nil {
			app.Logger.Fatal("couldn't load template data for "+name, "err", err)
		} else {
			data.StringMap["GeneratedAt"] = template.GeneratedAt.Format(time.UnixDate)
		}

		app.Logger.Debug("handling a links page", "data", &data)
		render.RenderTemplateWithData(writer, page, &data)
	}
}
