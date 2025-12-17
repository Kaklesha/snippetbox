package main

import (
	"net/http"
)

// the serverError helper writes a log entry at Error level (including the request
// method and URL as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.String()
		////added stack traces
		//trace = string(debug.Stack())
	)
	////Below commented Stack traces version logger
	//app.logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
	app.logger.Error(err.Error(), "method", method, "url", url)
}

// the clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)

}
