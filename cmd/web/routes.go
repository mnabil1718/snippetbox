package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) generateRoutes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequests, secureHeaders)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
