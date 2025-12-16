package main

import (
	"flag"
	"log"
	"net/http"
)

//UDP: Info about SERVE a SINGLE FILE did move to README read about it there

func main() {

	//Define a new command-line flag with the name 'addr', a default value of ":4000"
	//and some short help text explaning what the flag controls. The value of the
	//flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	//Importantly , we use the flag.Parse() function to parse the command-line flag.
	//This reads in the command-line flag value and assigns it to the addr
	//otherwise it will always contain the default value of ":4000". If any errors are
	//encountered during parsing the application will be terminated.
	flag.Parse()

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

	//The value returned from the flag.String() function is a pointer to the flag
	//value, not the value itself. So in this code, that means the addr variable
	//is actually a pointer, and we need to dereference it (i.e. prefix it with
	//the * symbol) before using it. Note that we're using the log.Print()
	//function to interpolate the address with the log message.
	log.Printf("starting server on %s", *addr)

	//And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	//go run ./cmd/web -addr=":9999"
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)

}
