package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies fot the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger *slog.Logger
}

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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		////Minimum log lvl: override lvl severity
		// Level: slog.LevelDebug,
		////Caller location: absolute path to file where invoke a slog. Info() and its line
		//AddSource: true,
	}))

	//Initialize a new instance of out application struct, containing the
	//dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
	}

	//ROUTING REST API to HANDLERS the SECTION
	logger.Info("starting server on ", slog.String("addr", *addr))
	//And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	//go run ./cmd/web -addr=":9999"
	err := http.ListenAndServe(*addr, app.routes())
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
