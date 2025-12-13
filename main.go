package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")     //need type before WriteHEader
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
	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))

}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))

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
