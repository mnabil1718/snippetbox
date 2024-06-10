package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) ServeError(writer http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLogger.Output(2, trace)
	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(writer http.ResponseWriter, statusCode int) {
	http.Error(writer, http.StatusText(statusCode), statusCode)
}

func (app *Application) NotFound(writer http.ResponseWriter) {
	http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
