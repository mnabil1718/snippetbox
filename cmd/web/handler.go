package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *Application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // root path restriction
		app.NotFound(writer)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		app.ServeError(writer, err)
		return
	}

	err = template.Execute(writer, nil)
	if err != nil {
		app.ServeError(writer, err)
	}
}
func (app *Application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.NotFound(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "Displaying snippet id: %d", id)
}
func (app *Application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
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

// doesnt need to use PascalCase, since its used only in the same package (main)
var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
