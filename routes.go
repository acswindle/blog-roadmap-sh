package main

import (
	"net/http"
	"strconv"
)

func (app application) HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := TemplateData{
		&app.articles,
		"Home",
		ArticleForm{},
	}
	app.ExecuteTemplate("home.tmpl.html", data, w)
}

func (app application) ArticleHandle(w http.ResponseWriter, r *http.Request) {
	idx, err := strconv.Atoi(r.PathValue("idx"))
	if err != nil {
		http.Error(w, "bad article value", http.StatusBadRequest)
	}
	if (idx >= len(app.articles)) || (idx < 0) {
		http.Error(w, "bad article value", http.StatusBadRequest)
	}
	data := TemplateData{
		&app.articles,
		"Home",
		ArticleForm{
			app.articles[idx],
			nil,
		},
	}
	app.ExecuteTemplate("article.tmpl.html", data, w)
}
