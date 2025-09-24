package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type ArticleForm struct {
	Article
	Errors []error
}

type TemplateData struct {
	Articles  *Articles
	PageTitle string
	ArticleForm
}

type templates map[string]*template.Template

func (app *application) RefreshTemplates() error {
	files, err := filepath.Glob(filepath.Join(app.templateDir, "*.tmpl.html"))
	if err != nil {
		return err
	}
	base := filepath.Join(app.templateDir, "base.tmpl.html")
	for _, file := range files {
		name := filepath.Base(file)
		if name == "base.tmpl.html" {
			continue
		}
		app.templates[name], err = template.ParseFiles(base, file)
		fmt.Println(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app application) ExecuteTemplate(name string, data TemplateData, w http.ResponseWriter) {
	tmpl, present := app.templates[name]
	if !present {
		http.Error(w, "failed loading home page", http.StatusInternalServerError)
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("error loading page %v", err)
	}
}
