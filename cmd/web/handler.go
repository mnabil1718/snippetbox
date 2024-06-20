package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mnabil1718/snippetbox/pkg/models"
)

func (app *Application) home(writer http.ResponseWriter, request *http.Request) {

	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServeError(writer, err)
		return
	}

	data := &TemplateData{Snippets: snippets}
	app.render(writer, "home.page.tmpl", data)
}

func (app *Application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get(":id"))

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
	data := &TemplateData{Snippet: snippet}
	app.render(writer, "show.page.tmpl", data)

}
func (app *Application) createSnippet(writer http.ResponseWriter, request *http.Request) {

	title := "Third entry"
	content := "Hello/nDarkness/My old friend."
	expiresInDays := "7"

	id, err := app.Snippets.Insert(title, content, expiresInDays)
	if err != nil {
		app.ServeError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *Application) createSnippetForm(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Create snippet form..."))
}

// doesnt need to use PascalCase, since its used only in the same package (main)
var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
