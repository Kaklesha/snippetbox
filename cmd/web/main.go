package main

import (
	"log"
	"net/http"
)

//UDP: Info about SERVE a SINGLE FILE did move to README read about it there

func main() {
	mux := http.NewServeMux()
	//UDP: Info about DISABLING dir listings move to README read about it there

	//Create a file server files out of the "./ui/static" directory.
	//Note that the  path given to the http.Dir function is relative to the project
	//directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "./static" . For matching paths, we strip the
	//"static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	//ROUTING REST API to HANDLERS the SECTION

	//This line using for example more detailes about homeD into handler.go
	mux.Handle("GET /homed/{$}", &homeD{})
	//EARLIER DONE BELOW
	//udp: delete {$}  for unity with main guide line below
	mux.HandleFunc("GET /", http.HandlerFunc(home))
	mux.HandleFunc("GET /snippet/view/{id}/", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
	log.Print("starting server on : 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
