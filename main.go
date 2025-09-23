package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Article struct {
	Title   string
	Date    string
	Content string
}

type Articles []Article

func (app *application) LoadArticles() {
	file, err := os.Open("./metadata.json")
	if err != nil {
		log.Fatalf("could not load metadata file: %v", err)
	}
	wrapper := map[string]Articles{}
	err = json.NewDecoder(file).Decode(&wrapper)
	if err != nil {
		log.Fatalf("could not decode metadata file %v", err)
	}
	var present bool
	app.articles, present = wrapper["articles"]
	if !present {
		log.Fatal("could not unwrap json file")
	}
	fmt.Printf("decoded json %v\n", app.articles)
}

type application struct {
	templateDir string
	templates   templates
	articles    Articles
}

func NewApp(templateDir string) (*application, error) {
	app := application{
		templateDir,
		make(templates),
		Articles{},
	}
	err := app.RefreshTemplates()
	if err != nil {
		return nil, err
	}
	app.LoadArticles()
	return &app, nil
}

func main() {
	app, err := NewApp("./templates/")
	if err != nil {
		log.Fatalf("could not load templates %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.HomeHandle)

	if err := http.ListenAndServe(
		":9090",
		mux,
	); err != nil {
		log.Fatalf("error starting server %v", err)
	}
}
