package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arammikayelyan/snippetbox/pkg/models"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always
		// be run in the event) of a panic as Go unwind the stack.
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response
				w.Header().Set("Connection", "close")
				// Call the app.serverError() method to return 500
				// Internal server error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retreive the userID from the session
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the details of the user from the db. If no matching
		// record is found, remove the (invalid) userID from the session
		// and call the next handler in the chain as normal.
		user, err := app.users.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		}
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Otherwise, we know that request is coming from the authenticated
		// user. We create a new copy of request with the user information
		// added to the context, and call the next handler in the chain *using
		// this new copy of the request*
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
