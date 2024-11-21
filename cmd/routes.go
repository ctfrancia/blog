package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /posts/{slug}", app.post)

	return mux
}
