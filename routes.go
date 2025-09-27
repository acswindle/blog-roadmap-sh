package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (app application) HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	articles, err := app.GetArticles()
	if err != nil {
		http.Error(w, "could not query articles", http.StatusInternalServerError)
		return
	}
	data := TemplateData{
		&articles,
		"Home",
	}
	app.ExecuteTemplate("home.tmpl.html", data, w)
}

func (app application) ArticleHandle(w http.ResponseWriter, r *http.Request) {
	idx, err := strconv.Atoi(r.PathValue("idx"))
	if err != nil {
		http.Error(w, "bad article value", http.StatusBadRequest)
	}
	article, err := app.GetArticle(idx)
	if err != nil {
		http.Error(w, "could not find article", http.StatusInternalServerError)
	}
	data := TemplateData{
		article,
		fmt.Sprintf("Article - %s", article.Title),
	}
	app.ExecuteTemplate("article.tmpl.html", data, w)
}

func (app application) CreateArticleHandle(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		ArticleForm{
			Article: Article{},
			Errors:  []error{},
		},
		"Create Article",
	}
	app.ExecuteTemplate("create.tmpl.html", data, w)
}

func (form *ArticleForm) NonEmptyString(input *string, key string, r *http.Request) {
	*input = r.Form.Get(key)
	*input = strings.TrimSpace(*input)
	if *input == "" {
		form.Errors = append(form.Errors, fmt.Errorf("must have a %s", key))
	}
}

func (app *application) PostArticleHandle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	form := ArticleForm{}
	form.NonEmptyString(&form.Title, "title", r)
	form.NonEmptyString(&form.Content, "content", r)
	if len(form.Errors) > 0 {
		fmt.Printf("error with field %v", form.Errors)
		app.ExecuteTemplate("create.tmpl.html", TemplateData{form, "Create Article"}, w)
		return
	}
	err := app.AddArticle(form.Article)
	if err != nil {
		http.Error(w, "error adding article", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
