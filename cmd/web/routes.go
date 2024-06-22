package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) generateRoutes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequests, secureHeaders)
	dynamicMiddleware := alice.New(app.Session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/ping", http.HandlerFunc(ping))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createSnippet))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createSnippetForm))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signup))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.login))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.logout))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
