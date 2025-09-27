package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type Article struct {
	ID      int
	Title   string
	Date    time.Time
	Content string
}

type Articles []Article

type application struct {
	templateDir string
	templates   templates
	db          *sql.DB
}

func NewApp(templateDir string) (*application, error) {
	db := NewDB()
	app := application{
		templateDir,
		make(templates),
		db,
	}
	err := app.RefreshTemplates()
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func main() {
	app, err := NewApp("./templates/")
	defer app.db.Close()
	if err != nil {
		log.Fatalf("could not load templates %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.HomeHandle)
	mux.HandleFunc("GET /article/create", app.CreateArticleHandle)
	mux.HandleFunc("POST /article/create", app.PostArticleHandle)
	mux.HandleFunc("GET /article/{idx}", app.ArticleHandle)

	if err := http.ListenAndServe(
		":9090",
		mux,
	); err != nil {
		log.Fatalf("error starting server %v", err)
	}
}
