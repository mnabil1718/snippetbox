package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello from the other side"))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	fmt.Println("Starting server on port:8080...")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)

}
