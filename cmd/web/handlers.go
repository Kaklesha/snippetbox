package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.kira.net/internal/models"
)

// UDP: Changed the signature of the home handler so it is defined as a method against
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	snippets, err := app.snippets.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//!IMPORTANT: This a little deviation from the author's implementation
	data, err := NewTemplateData(&snippets)
	if err != nil {
		app.serverError(w, r, err)
	}
	app.render(w, r, http.StatusOK, "home.html", *data)

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
	data, err := NewTemplateDatum(&snippet)
	if err != nil {
		app.serverError(w, r, err)
	}
	app.render(w, r, http.StatusOK, "view.html", *data)

}

// Add a snippetCreate handler function.
// Changed the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"name":"Display a form for creating a new snippet..."}`))

}

// Changed the signature of the snippetView handler so it is defined as a method
// against *application
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//use the method to send a 201 status code?

	//Create some valiables holding dummy data. We'll remove these later on
	//during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
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
