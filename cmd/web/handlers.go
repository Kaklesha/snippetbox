package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"snippetbox.kira.net/internal/models"
)

// UDP: Changed the signature of the home handler so it is defined as a method against
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	////Manipulating the header map
	//Set a new cache-control header . If an existing "Cache-Control" header exists
	// it will be overwritten.
	w.Header().Set("Cache-Control", "public, max-age=31536000 ")

	//In contrast, the Add() method appends a new "Cache-Control" header and can
	// be called multiple times.
	w.Header().Add("Cache-Control", "public")
	w.Header().Add("Cache-Control", "max-age=31536000")

	//Delete all values for the "Cacha-Control" header.
	w.Header().Del("Cache-Control")

	//Retrieve the first value for the "Cache-Control" header.
	w.Header().Get("Cache-Control")
	//Retrieve a slice of all values for the "Cache-Control" header.
	w.Header().Values("Cache-Control")
	//below comment oldCode
	//w.Header().Add("Server", "Go")     //need type before WriteHEader
	w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	w.WriteHeader(http.StatusAccepted) //need type defore Write
	////w.Write([]byte("HEllo from snippetbox"))

	// Initialize a slice containing the paths to the two files. It's important
	// to note that file file containing our base template must be the *first*
	//file in the slice

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	//Use the template.ParseFiles() function to read the template file into a
	//template set. If there's an error , we log the detailed error message, use
	//the http.Error() function to send an Internal Server Error response to the
	//  user , and then return from the handler so no subsequent code is executed.

	ts, err := template.ParseFiles(files...)
	if err != nil {
		//Because the home handler is now a method against the application
		//struct it can access its fields, including the structured logger.
		//We'll use this to create a log entry at Error level containing the error
		//message, also including the request method and URL as attributes to
		//assist with debugging.
		//!UDP in here below implemented r.URL.String() that returns the full URL stirng
		// Against r.URL.RequestURL that undefined (type *url.URL has no field or method RequestURL)
		app.serverError(w, r, err)
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//(DEPRECATE: we added ExecuteTemplate(_)) Then we use the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	//leave as nil

	//Use the ExecuteTemplate() method to write the content of the "base"
	//template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		app.serverError(w, r, err)
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", snippet)

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
