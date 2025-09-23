package main

import "net/http"

func (app application) HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := TemplateData{
		&app.articles,
		"Home",
		nil,
	}
	app.ExecuteTemplate("home.tmpl.html", data, w)
}
