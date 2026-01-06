package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	"snippetbox.kira.net/internal/models"

	"github.com/go-playground/form/v4"

	"github.com/alexedwards/scs/mysqlstore" // New import
	"github.com/alexedwards/scs/v2"         // New import

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies fot the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
// UDP
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	//Define a new command-line  flag for the MySql DSN string.
	dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	//use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		////Minimum log lvl: override lvl severity
		// Level: slog.LevelDebug,
		////Caller location: absolute path to file where invoke a slog. Info() and its line
		//AddSource: true,
	}))

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	//initialize a decoder instance
	formDecoder := form.NewDecoder()

	//Use the scs.New() function to initialize a new session manager. Then we
	//configure it to use our MySQL database as the  session store , and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created.)
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	//By default set SameSite=Lax instead of SameSite=Struct as below
	//	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	//Initialize a new instance of out application struct, containing the
	//dependencies (for now, just the structured logger).
	//UDP Add it tothe application dependencies.

	//Initialize a models/UserModel instance and dd it to the application
	//dependencies.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	//Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	//is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		MinVersion: tls.VersionTLS13,
	}

	//Initialize a new http.Server struct. We set the Addr and Handler fields so
	//that the server usess the same network address and routes as before.
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at Error level, and assign it to the ErrorLog field. If
		// you would prefer to log the server errors at Warn level instead, you
		// could pass slog.LevelWarn as the final parameter.
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		TLSConfig: tlsConfig,
		//Add idle , read and write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//ROUTING REST API to HANDLERS the SECTION
	logger.Info("starting server", slog.String("addr", *addr))
	//And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	//go run ./cmd/web -addr=":9999"
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	logger.Error(err.Error())
	os.Exit(1)
}

// the openDB function wraps sql.Open() and retuns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
