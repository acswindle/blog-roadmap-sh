package main

import "net/http"

func (app application) HomeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	app.ExecuteTemplate("home.tmpl.html", nil, w)
}
