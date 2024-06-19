package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mnabil1718/snippetbox/pkg/models/postgresql"
)

type Application struct {
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	Snippets      *postgresql.SnippetModel
	TemplateCache map[string]*template.Template
}

func OpenDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	// create log file
	file := CreateFile("./logs/", "./logs/application.log")

	addr := flag.String("addr", "localhost:8080", "HTTP Server Address")
	dsn := flag.String("dsn", "postgres://mnabil:Cucibaju123@localhost:5432/snippetbox", "PostgreSQL Connection String")
	flag.Parse() // add any custom flag before parsing

	db, err := OpenDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cache, err := newTemplateCache("./ui/html")
	if err != nil {
		log.Fatal(err)
	}
	// entry point for dependency injection
	app := &Application{
		InfoLogger:    log.New(file, "INFO \t", log.Ldate|log.Ltime),
		ErrorLogger:   log.New(file, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile),
		Snippets:      &postgresql.SnippetModel{DB: db},
		TemplateCache: cache,
	}

	server := http.Server{
		Addr:     *addr,
		Handler:  app.generateRoutes(),
		ErrorLog: app.ErrorLogger, // only for HTTP errors
	}

	app.InfoLogger.Printf("Starting server on %s...", *addr)
	err = server.ListenAndServe()
	app.ErrorLogger.Fatal(err)
}
