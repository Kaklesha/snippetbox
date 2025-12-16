package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
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

	//use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
	//UDP: log.Print to logger. Info()
	//Use the Info() method to log the starting server message at Info severity
	//(along) with the listen address as an attribute).
	//UDP: use slog.String() to create atttibutes is more verbose
	//but safer in sense that it reduces the risk of bugs in your application.
	logger.Info("starting server on ", slog.String("addr", *addr))

	//And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	//go run ./cmd/web -addr=":9999"
	err := http.ListenAndServe(*addr, mux)
	//And we also use the Error() method to log any  error message returned by
	//http.ListenAndServe() at Error severity (with no additional attributes),
	//and then call os.Exit(1) to terminate the application with exit code 1.
	//IMPORT!! There is no structured logging equivalent to the log.Fatal() function
	// that we can use to handle an error returned by http.ListenAndServe().Instead, the
	// closest we can get is logging a message at the Error severity level and then manually
	// calling os.Exit(1) to terminate the application with the exit code 1, like we are in the
	// code above.
	logger.Error(err.Error())
	os.Exit(1)
}
