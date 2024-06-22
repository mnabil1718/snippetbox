package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/mnabil1718/snippetbox/pkg/models"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-XSS-Protection", "1; mode=block")
		writer.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(writer, request)
	})
}

func (app *Application) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.Session.Put(r, "redirectPathAfterLogin", r.URL.Path)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// remember this whole middleware is just to fix isAuthenticated helper function
// because previously it just checks if isAuthenticatedUserID key exists in session or not
// This middleware only do checks to the database if isAuthenticatedUserID is present.
func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		exists := app.Session.Exists(r, authenticatedUserIDSessionKey)

		// NON-AUTHENTICATED | NON-VERIFIED
		// if session authID not exists, then pass it to requireAuth middleware as is
		// or maybe this requests go to unauthenticated routes, which is fine.
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// AUTHENTICATED | NON-VERIFIED
		// if authID key exists, we can't trust it fully yet.
		// we need to validate to DB if this ID for this user is correct
		// if its not, we are going to treat them as unauthenticated
		user, err := app.Users.Get(app.Session.GetInt(r, authenticatedUserIDSessionKey))
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.Session.Remove(r, authenticatedUserIDSessionKey)
				next.ServeHTTP(w, r)
				return
			}

			app.ServerError(w, err)
			return
		}

		if !user.Active {
			app.Session.Remove(r, authenticatedUserIDSessionKey)
			next.ServeHTTP(w, r)
			return
		}

		// AUTHENTICATED | VERIFIED
		// all checks pass, we create child context with value
		// representing if user is authenticated and verified
		ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
