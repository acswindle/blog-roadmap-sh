package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type templates map[string]*template.Template

func (app *application) RefreshTemplates() error {
	files, err := filepath.Glob(filepath.Join(app.templateDir, "*.tmpl.html"))
	if err != nil {
		return err
	}
	for _, file := range files {
		name := filepath.Base(file)
		app.templates[name], err = template.ParseFiles(file)
		fmt.Println(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app application) ExecuteTemplate(name string, data any, w http.ResponseWriter) {
	tmpl, present := app.templates[name]
	if !present {
		http.Error(w, "failed loading home page", http.StatusInternalServerError)
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "failed loading home page", http.StatusInternalServerError)
	}
}
