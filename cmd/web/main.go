package main

import (
	"log"
	"net/http"
)

// // Sometimes you might want to serve a single
// // file from within a handler. For this there’s the
// // http.ServeFile() function, which you can use like so:
// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// http.ServeFile(w, r, "./ui/static/file.zip")
// http.ServerFile() does not automatically sanitize the file path.
// //you must sanitize the input with "filepath.Clean()"before using it.
// }
func main() {
	mux := http.NewServeMux()
	//FOR DISABLING DIRECTORY LISTINGS thinking
	// If you want to disable directory listings there are a few different approaches you can take.
	// The simplest way? Add a blank index.html file to the specific directory that you want to
	// disable listings for. This will then be served instead of the directory listing, and the user will
	// get a 200 OK response with no body. If you want to do this for all directories under
	// ./ui/static you can use the command:
	// $ find./ui/static -type d -exec touch {}/index.html \;
	// A more complicated (but arguably better) solution is to create a custom implementation of
	// http.FileSystem, and have it return an os.ErrNotExist error for any directories. A full
	// explanation and sample code can be found in this blog post.

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
	mux.HandleFunc("GET /", home)
	mux.HandleFunc("GET /snippet/view/{id}/", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
	log.Print("starting server on : 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
