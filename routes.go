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

func (app *application) GetRequestArticle(w http.ResponseWriter, r *http.Request) (Article, bool) {
	idx, err := strconv.Atoi(r.PathValue("idx"))
	if err != nil {
		http.Error(w, "bad article value", http.StatusBadRequest)
		return Article{}, false
	}
	article, err := app.GetArticle(idx)
	if err != nil {
		http.Error(w, "could not find article", http.StatusInternalServerError)
		return Article{}, false
	}
	return article, true
}

func (app application) ArticleHandle(w http.ResponseWriter, r *http.Request) {
	article, ok := app.GetRequestArticle(w, r)
	if !ok {
		return
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

func (app application) EditArticleHandle(w http.ResponseWriter, r *http.Request) {
	article, ok := app.GetRequestArticle(w, r)
	if !ok {
		return
	}
	data := TemplateData{
		ArticleForm{
			Article: article,
			Errors:  []error{},
		},
		"Edit Article",
	}
	app.ExecuteTemplate("edit.tmpl.html", data, w)
}

func (form *ArticleForm) NonEmptyString(input *string, key string, r *http.Request) {
	*input = r.Form.Get(key)
	*input = strings.TrimSpace(*input)
	if *input == "" {
		form.Errors = append(form.Errors, fmt.Errorf("must have a %s", key))
	}
}

func (app *application) CheckArticleForm(w http.ResponseWriter, r *http.Request, webpage string) (*ArticleForm, bool) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return &ArticleForm{}, false
	}
	form := ArticleForm{}
	form.NonEmptyString(&form.Title, "title", r)
	form.NonEmptyString(&form.Content, "content", r)
	if len(form.Errors) > 0 {
		app.ExecuteTemplate(webpage, TemplateData{form, "Create Article"}, w)
		return &ArticleForm{}, false
	}
	return &form, true
}

func (app *application) PostArticleHandle(w http.ResponseWriter, r *http.Request) {
	form, ok := app.CheckArticleForm(w, r, "create.tmpl.html")
	if !ok {
		return
	}
	err := app.AddArticle(form.Article)
	if err != nil {
		http.Error(w, "error adding article", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) EditArticlePostHandle(w http.ResponseWriter, r *http.Request) {
	idx, err := strconv.Atoi(r.PathValue("idx"))
	if err != nil {
		http.Error(w, "invalid article id in path", http.StatusBadRequest)
		return
	}
	form, ok := app.CheckArticleForm(w, r, "edit.tmpl.html")
	if !ok {
		return
	}
	form.ID = idx
	err = app.UpdateArticle(form.Article)
	if err != nil {
		fmt.Printf("error editing artilce %v", err)
		http.Error(w, "error editing article", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
