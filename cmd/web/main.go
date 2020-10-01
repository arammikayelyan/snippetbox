package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/arammikayelyan/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide
// dependencies.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Parsing the runtime configs for the app.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass123!@/snippetbox?parseTime=true", "MySQL dsn")
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

	// Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
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
