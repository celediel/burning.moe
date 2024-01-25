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
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.Logger.Info("REQUEST", "url", r.URL, "ip", r.RemoteAddr, "useragent", r.UserAgent())
			next.ServeHTTP(w, r)
		})
	},
	middleware.Recoverer,
}
