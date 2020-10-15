package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the
// errorLog, then sends a generic 500 Internal Server Error response to
// the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// The clientError helper sends a specific status code and corresponding
// description to the user. We'll use this later when there's a problem
// the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement NotFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found
// response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Create an addDefaultData helper. This takes a pointer to templateData struct,
// adds the current year to the CurrentYear field, and then returns the pointer.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	// Add the flash message to the templateData, if one exists.
	td.Flash = app.session.PopString(r, "flash")
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the apropriate template set from the cache based on the page name.
	// If no entry exists in the cache with the provided name, call the serverError
	// method.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// authenticatedUser returns the ID of the user from the session,
// or zero if the request is from the an unauthenticated user.
func (app *application) authenticatedUser(r *http.Request) int {
	return app.session.GetInt(r, "userID")
}
