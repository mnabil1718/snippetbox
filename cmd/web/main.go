package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", snippetCreate)

	fmt.Println("Starting server on port:8080...")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)

}
