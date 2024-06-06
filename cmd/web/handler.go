package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // root path restriction
		http.NotFound(writer, request)
		return
	}

	template, err := template.ParseFiles("./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
		return
	}

	err = template.Execute(writer, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
	}
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
