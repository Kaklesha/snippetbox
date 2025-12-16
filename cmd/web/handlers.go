package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Additional information: example using
// Instead of func named "home" we take empty a type named "homeD"
// for demo corresponding mwthod ServerHTTP
type homeD struct{}

func (h *homeD) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is my homeD page"))
}

func home(w http.ResponseWriter, r *http.Request) {
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
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// add a snippetView  handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	//Extract the value of the id wildcard from the request using r.PathValue()
	//and try to convert it to an integer using the strconv.Atoi() function. If
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"name":"Display a form for creating a new snippet..."}`))

}
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//use the method to send a 201 status code
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
	//
}
