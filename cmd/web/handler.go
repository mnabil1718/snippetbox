package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mnabil1718/snippetbox/pkg/models"
)

func (app *Application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" { // root path restriction
		app.NotFound(writer)
		return
	}

	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServeError(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	for _, snippet := range snippets {
		fmt.Fprintf(writer, "%v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// template, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.ServeError(writer, err)
	// 	return
	// }

	// err = template.Execute(writer, nil)
	// if err != nil {
	// 	app.ServeError(writer, err)
	// }
}
func (app *Application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.NotFound(writer)
		return
	}

	snippet, err := app.Snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(writer)
		} else {
			app.ServeError(writer, err)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "%v", snippet)
}
func (app *Application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allowed", http.MethodPost)
		writer.Header().Set("Content-Type", "application/json")
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := app.Snippets.Insert("Third entry", "Hello/nDarkness/My old friend.", "7")
	if err != nil {
		app.ServeError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

// doesnt need to use PascalCase, since its used only in the same package (main)
var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
