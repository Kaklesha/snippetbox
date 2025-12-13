package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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
	w.Write([]byte("HEllo from snippetbox"))

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
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}/", snippetView)
	//mux.HandleFunc("/snippet/view/remainder/{path...}", snippetRemanderView)
	//example remainer wildcards is msg :
	//http://localhost:4000/snippet/view/remainder/sdf/sdfs/sdf/2/34/2
	// Display a specific snippet with ID %d...sdf/sdfs/sdf/2/34/2
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
	log.Print("starting server on : 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
