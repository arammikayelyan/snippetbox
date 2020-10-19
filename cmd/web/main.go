package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/arammikayelyan/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

// application struct holds the application-wide dependencies.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {
	// Parsing the runtime configs for the app.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass123!@/snippetbox?parseTime=true", "MySQL dsn")

	// Define a new command-line flag for the session secret
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret session")
	flag.Parse()

	// Creating a logger for writing information and error
	// messages to stdout (os.Stdout) and stderr (os.Stderr).
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// Pass openDB() dsn from command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Close connection pool before the main() function exits.
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	// Initialize a custom server with address, custom logger, and handler.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Running the HTTP server.
	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
