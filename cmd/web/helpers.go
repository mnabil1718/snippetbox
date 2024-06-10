package main

import (
	"fmt"
	"net/http"
	"os"
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

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateFile(parentDirPath string, filePath string) *os.File {

	if parentDirPath != "" {
		err := os.MkdirAll(parentDirPath, 0755)
		PanicIfError(err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	PanicIfError(err)

	return file
}
