package main

import (
	"flag"
	"log"
	"net/http"
)

type Application struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func main() {

	// create log file
	file := CreateFile("./logs/", "./logs/application.log")

	// entry point for dependency injection
	app := &Application{
		InfoLogger:  log.New(file, "INFO \t", log.Ldate|log.Ltime),
		ErrorLogger: log.New(file, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	addr := flag.String("addr", "localhost:8080", "HTTP Server Address")
	flag.Parse()

	server := http.Server{
		Addr:     *addr,
		Handler:  app.generateRoutes(),
		ErrorLog: app.ErrorLogger, // only for HTTP errors
	}

	app.InfoLogger.Printf("Starting server on %s...", *addr)
	err := server.ListenAndServe()
	app.ErrorLogger.Fatal(err)
}
