package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	// create necessary loggers
	// create logs dir first
	err := os.MkdirAll("./logs/", 0755)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("./logs/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	infoLog := log.New(file, "INFO \t", log.Ldate|log.Ltime)
	errorLog := log.New(file, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", "localhost:8080", "HTTP Server Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	server := http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog, // only for HTTP errors
	}

	infoLog.Printf("Starting server on %s...", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
