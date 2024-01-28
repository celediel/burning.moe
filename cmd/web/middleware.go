package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware is a slice of Middleware (aka func(n http.Handler) http.Handler {})
var Middleware []func(next http.Handler) http.Handler = []func(next http.Handler) http.Handler{
	// chi's recommended list
	middleware.RequestID,
	middleware.RealIP,
	// plus custom request logger
	func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			app.Logger.Info("REQUEST", "url", request.URL, "ip", request.RemoteAddr, "useragent", request.UserAgent())
			next.ServeHTTP(writer, request)
		})
	},
	middleware.Recoverer,
}
