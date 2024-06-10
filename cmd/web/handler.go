package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *Application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // root path restriction
		http.NotFound(writer, request)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLogger.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
		return
	}

	err = template.Execute(writer, nil)
	if err != nil {
		app.ErrorLogger.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
	}
}
func (app *Application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("Display spewcific snippet"))
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
