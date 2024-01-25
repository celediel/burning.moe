package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware is a slice of Middleware (aka func(n http.Handler) http.Handler {})
var Middleware []func(next http.Handler) http.Handler = []func(next http.Handler) http.Handler{
	middleware.Recoverer,
}
