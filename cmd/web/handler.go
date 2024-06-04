package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // root path restriction
		http.NotFound(writer, request)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("Hello from the other side"))
}
func showSnippet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("Display spewcific snippet"))
}
func snippetCreate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allowed", http.MethodPost)
		writer.Header().Set("Content-Type", "application/json")
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil || id < 1 {
		writer.Header().Set("Content-Type", "application/json")
		http.Error(writer, "ID Not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "Created snippet id: %d", id)
}
