package main

import (
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/mnabil1718/snippetbox/pkg/forms"
	"github.com/mnabil1718/snippetbox/pkg/models"
)

type TemplateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            *forms.Form
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func newlineToBreak(s string) template.HTML {
	str := strings.Replace(s, "\\n", "<br>", -1)
	return template.HTML(str)
}

var funcMap = template.FuncMap{
	"humanDate":      humanDate,
	"newlineToBreak": newlineToBreak,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(funcMap).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet

	}

	return cache, nil
}
