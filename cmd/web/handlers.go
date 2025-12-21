package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"snippetbox.kira.net/internal/models"
)

// UDP: Changed the signature of the home handler so it is defined as a method against
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.html", data)

}

// add a snippetView  handler function.
// Changed the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	//Extract the value of the id wildcard from the request using r.PathValue()
	//and try to convert it to an integer using the strconv.Atoi() function. If
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Use the SnippetModel's Get() method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//!IMPORTANT: This a little deviation from the author's implementation
	// data, err := NewTemplateDatum(&snippet)

	data := app.newTemplateData(r)
	data.Snippet = &snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

// Add a snippetCreate handler function.
// Changed the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.html", data)
}

// Changed the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return

	}
	// Initialize a map to hold any validation errors for the form fields.
	fieldErrors := make(map[string]string)
	// Check that the title value is not blank and is not more than 100
	// characters long. If it fails either of those checks, add a message to the
	// errors map using the field name as the key.
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	// Check that the Content value isn't blank.
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}
	// Check the expires value matches one of the permitted values (1, 7 or
	// 365).
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}
	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
func (app *application) snippetTransact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := app.snippets.ExampleTransaction()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Write([]byte(`{"transaction":"Display about successful..."}`))

}
