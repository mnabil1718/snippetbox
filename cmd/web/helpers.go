package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"
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

func (app *Application) render(writer http.ResponseWriter, request *http.Request, name string, templateData *TemplateData) {
	templateSet, ok := app.TemplateCache[name]
	if !ok {
		app.ServeError(writer, fmt.Errorf("template name %s doesn't exists", name))
	}

	bufferPointer := new(bytes.Buffer)

	err := templateSet.Execute(bufferPointer, app.addDefaultData(request, templateData))
	if err != nil {
		app.ServeError(writer, err)
		return
	}

	bufferPointer.WriteTo(writer)
}

func (app *Application) addDefaultData(request *http.Request, templateData *TemplateData) *TemplateData {
	flash := app.Session.PopString(request, "flash")
	if templateData == nil {
		return &TemplateData{CurrentYear: time.Now().Year(), Flash: flash}
	}
	templateData.Flash = flash
	templateData.CurrentYear = time.Now().Year()
	return templateData
}
