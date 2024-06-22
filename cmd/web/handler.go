package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mnabil1718/snippetbox/pkg/forms"
	"github.com/mnabil1718/snippetbox/pkg/models"
)

func ping(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("OK"))
}

func (app *Application) home(writer http.ResponseWriter, request *http.Request) {

	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServerError(writer, err)
		return
	}

	data := &TemplateData{Snippets: snippets}
	app.render(writer, request, "home.page.tmpl", data)
}

func (app *Application) about(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "about.page.tmpl", nil)
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
			app.ServerError(writer, err)
		}
		return
	}

	data := &TemplateData{Snippet: snippet}
	app.render(writer, request, "show.page.tmpl", data)
}

func (app *Application) createSnippet(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		app.ClientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)

	form.Required("title", "content", "expires")
	form.MaxLength(100, "title")
	form.PermittedValues("expires", "1", "7", "365")

	if !form.Valid() {
		app.render(writer, request, "create.page.tmpl", &TemplateData{Form: form})
		return
	}

	id, err := app.Snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.ServerError(writer, err)
		return
	}

	app.Session.Put(request, "flash", "Snippet created sucessfully")
	http.Redirect(writer, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *Application) createSnippetForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "create.page.tmpl", &TemplateData{Form: forms.New(nil)})
}

func (app *Application) signupForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "signup.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *Application) signup(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.ClientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)

	form.Required("name", "email", "password")
	form.MaxLength(255, "name", "email")
	form.MatchesPattern(forms.EmailRX, "email")
	form.MinLength(10, "password")

	if !form.Valid() {
		app.render(writer, request, "signup.page.tmpl", &TemplateData{Form: form})
		return
	}

	err = app.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "email is already in use.")
			app.render(writer, request, "signup.page.tmpl", &TemplateData{Form: form})
			return
		}

		app.ServerError(writer, err)
		return
	}

	app.Session.Put(request, "flash", "User registered sucessfully.")
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *Application) loginForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "login.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *Application) login(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.ClientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)

	form.Required("email", "password")
	form.MatchesPattern(forms.EmailRX, "email")

	if !form.Valid() {
		app.render(writer, request, "login.page.tmpl", &TemplateData{Form: form})
		return
	}

	id, err := app.Users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect.")
			app.render(writer, request, "login.page.tmpl", &TemplateData{Form: form})
			return
		}
		app.ServerError(writer, err)
		return
	}

	app.Session.Put(request, "authenticatedUserID", id)
	app.Session.Put(request, "flash", "Login Successful.")

	path := app.Session.PopString(request, "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(writer, request, path, http.StatusSeeOther)
		return
	}
	http.Redirect(writer, request, "/snippet/create", http.StatusSeeOther)
}

func (app *Application) logout(writer http.ResponseWriter, request *http.Request) {
	app.Session.Remove(request, "authenticatedUserID")
	app.Session.Put(request, "flash", "Logout Successful.")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func (app *Application) profile(writer http.ResponseWriter, request *http.Request) {
	id := app.Session.GetInt(request, authenticatedUserIDSessionKey)

	user, err := app.Users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.Session.Remove(request, authenticatedUserIDSessionKey)
			app.NotFound(writer)
			return
		}

		app.ServerError(writer, err)
		return
	}

	app.render(writer, request, "profile.page.tmpl", &TemplateData{User: user})
}

func (app *Application) changePasswordForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "password.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *Application) changePassword(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.ClientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)

	form.Required("currentPassword", "newPassword", "newPasswordConfirmation")
	form.MinLength(10, "currentPassword", "newPassword", "newPasswordConfirmation")
	form.IsEqual("newPassword", "newPasswordConfirmation")

	if !form.Valid() {
		app.render(writer, request, "password.page.tmpl", &TemplateData{Form: form})
		return
	}

	err = app.Users.ChangePassword(app.Session.GetInt(request, authenticatedUserIDSessionKey), form.Get("currentPassword"), form.Get("newPassword"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("currentPassword", "invalid password.")
			app.render(writer, request, "password.page.tmpl", &TemplateData{Form: form})
			return
		}

		app.ServerError(writer, err)
		return
	}

	app.Session.Put(request, "flash", "Password changed sucessfully.")
	http.Redirect(writer, request, "/user/profile", http.StatusSeeOther)
}

// doesnt need to use PascalCase, since its used only in the same package (main)
var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
